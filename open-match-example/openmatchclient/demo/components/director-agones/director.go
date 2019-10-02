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

package director

import (
	"io/ioutil"
	"log"
	"context"
	"fmt"
	"encoding/json"
	"io"
	"math/rand"
	"crypto/tls"
	"net/http"
	"time"
	"github.com/sirupsen/logrus"
	"open-match-example/openmatchclient/demo/components"
	"open-match-example/openmatchclient/demo/internal/rpc"
	"open-match.dev/open-match/pkg/pb"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "examples.demo",
	})
)

func Run(ds *components.DemoShared) {
	for !isContextDone(ds.Ctx) {
		run(ds)
	}
}

func isContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

type status struct {
	Status        string
	LatestMatches []*pb.Match
}

func run(ds *components.DemoShared) {
	defer func() {
		r := recover()
		if r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}

			ds.Update(status{Status: fmt.Sprintf("Encountered error: %s", err.Error())})
			time.Sleep(time.Second * 10)
		}
	}()

	s := status{}

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Connecting to backend"
	ds.Update(s)

	conn, err := rpc.GRPCClientFromConfig(ds.Cfg, "api.backend")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	be := pb.NewBackendClient(conn)

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Match Match: Sending Request"
	ds.Update(s)

	var matches []*pb.Match
	{
		req := &pb.FetchMatchesRequest{
			Config: &pb.FunctionConfig{
				Host: ds.Cfg.GetString("api.functions.hostname"),
				Port: int32(ds.Cfg.GetInt("api.functions.grpcport")),
				Type: pb.FunctionConfig_GRPC,
			},
			Profiles: []*pb.MatchProfile{
				{
					Name: "levels",
					Pools: []*pb.Pool{
						{
							Name: "Level1",
							StringEqualsFilters: []*pb.StringEqualsFilter{
								{
									Attribute: "player.level",
								    Value: "Level1",
								},
							},
						},
						{
							Name: "Level2",
							StringEqualsFilters: []*pb.StringEqualsFilter{
								{
									Attribute: "player.level",
								    Value: "Level2",
								},
							},
						},
						{
							Name: "Level3",
							StringEqualsFilters: []*pb.StringEqualsFilter{
								{
									Attribute: "player.level",
								    Value: "Level3",
								},
							},
						},
						{
							Name: "Level4",
							StringEqualsFilters: []*pb.StringEqualsFilter{
								{
									Attribute: "player.level",
								    Value: "Level4",
								},
							},
						},	
						{
							Name: "Level5",
							StringEqualsFilters: []*pb.StringEqualsFilter{
								{
									Attribute: "player.level",
								    Value: "Level5",
								},
							},
						},												
					},
				},
			},
		}

		stream, err := be.FetchMatches(ds.Ctx, req)
		if err != nil {
			logger.Info(fmt.Sprintf("Encountered error: %s", err.Error()))
			panic(err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Info(fmt.Sprintf("Encountered error: %s", err.Error()))
				panic(err)
			}
			matches = append(matches, resp.GetMatch())
		}
	}
	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Matches Found"
	s.LatestMatches = matches
	ds.Update(s)

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Assigning Players"
	ds.Update(s)

	for _, match := range matches {
		ids := []string{}

		for _, t := range match.Tickets {
			ids = append(ids, t.Id)
		}
        address, port := allocate()
		req := &pb.AssignTicketsRequest{
			TicketIds: ids,
			Assignment: &pb.Assignment{
				Connection: fmt.Sprintf("%s:%s", address, port),
			},
		}

		resp, err := be.AssignTickets(ds.Ctx, req)
		if err != nil {
			panic(err)
		}

		_ = resp
	}

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Sleeping"
	ds.Update(s)

	time.Sleep(time.Second * 5)
}

func allocate() (string, string) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}		
    var username string = "v1GameClientKey"
    var passwd string = "EAEC945C371B2EC361DE399C2F11E"
    req, err := http.NewRequest("GET", "https://allocator.agones.svc.cluster.local/address", nil)
    req.SetBasicAuth(username, passwd)
    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }
    bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
    var allocation map[string]interface{}
    json.Unmarshal([]byte(s), &allocation)
    status := allocation["status"].(map[string]interface{})
    address := status["address"].(string)
    portsmap := status["ports"]
    portslist := portsmap.([]interface {})
    gameportmap := (portslist[0].(map[string]interface{}))
    gameport := fmt.Sprintf("%.0f", gameportmap["port"])
    return address, gameport
}
