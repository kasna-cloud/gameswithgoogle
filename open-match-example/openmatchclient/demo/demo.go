// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package demo contains the core startup code for running a demo.
package demo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"open-match-example/openmatchclient/demo/bytesub"
	"open-match-example/openmatchclient/demo/components"
	"open-match-example/openmatchclient/demo/updater"
	"open-match-example/openmatchclient/demo/internal/config"
	"open-match-example/openmatchclient/demo/internal/logging"
	"open-match-example/openmatchclient/demo/internal/telemetry"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "examples.demo",
	})
)

// Run starts the provided components, and hosts a webserver for observing the
// output of those components.
func Run(comps map[string]func(*components.DemoShared)) {
	cfg, err := config.Read()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot read configuration.")
	}
	logging.ConfigureLogging(cfg)

	logger.Info("Initializing Server")

	fileServe := http.FileServer(http.Dir("/app/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServe))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fileServe.ServeHTTP(w, r)
	})

	http.Handle(telemetry.HealthCheckEndpoint, telemetry.NewAlwaysReadyHealthCheck())

	bs := bytesub.New()
	u := updater.New(context.Background(), func(b []byte) {
		var out bytes.Buffer
		err := json.Indent(&out, b, "", "  ")
		if err == nil {
			bs.AnnounceLatest(out.Bytes())
		} else {
			bs.AnnounceLatest(b)
		}
	})

	http.Handle("/connect", websocket.Handler(func(ws *websocket.Conn) {
		bs.Subscribe(ws.Request().Context(), ws)
	}))

	logger.Info("Starting Server")

	for name, f := range comps {
		go f(&components.DemoShared{
			Ctx:    context.Background(),
			Cfg:    cfg,
			Update: u.ForField(name),
		})
	}

	address := fmt.Sprintf(":%d", cfg.GetInt("api.demo.httpport"))
	err = http.ListenAndServe(address, nil)
	logger.WithError(err).Warning("HTTP server closed.")
}
