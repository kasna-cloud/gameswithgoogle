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

package clients

import (
	"context"
	"math/rand"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"open-match-example/openmatchclient/demo/components"
	"open-match-example/openmatchclient/demo/updater"
	"open-match-example/openmatchclient/demo/internal/config"
	"open-match-example/openmatchclient/demo/internal/rpc"
	"open-match.dev/open-match/pkg/pb"
	"open-match.dev/open-match/pkg/structs"
)


var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "examples.demo",
	})
)

func Run(ds *components.DemoShared) {
	// every minute new players
	logger.Info("client function run")
	levels := [5]string{"Level1", "Level2", "Level3", "Level4", "Level5"}
	u := updater.NewNested(ds.Ctx, ds.Update)
	for i := 0; i < 100; i++ {
		logger.Info(fmt.Sprintf("Creating Batch: %d", i))
		runOn(func (ds *components.DemoShared) {
			for j := 0; j < 2; j++ {
				name := fmt.Sprintf("batch_%d_fakeplayer_%d", i, j)
				exp := random(1, 5000)
				level := levels[rand.Intn(len(levels))]
				go func() {
					for !isContextDone(ds.Ctx) {
						runScenario(ds.Ctx, ds.Cfg, name, exp, level, u.ForField(name))
					}
				}()
			}
		}, ds )
		time.Sleep(15 * time.Second)
	}
}

func runOn(f func(ds *components.DemoShared), ds *components.DemoShared) {
	f(ds)
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
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
	Status     string
	Assignment *pb.Assignment
}

func runScenario(ctx context.Context, cfg config.View, name string, exp int, level string, update updater.SetFunc) {
	defer func() {
		r := recover()
		if r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}

			update(status{Status: fmt.Sprintf("Encountered error: %s", err.Error())})
			time.Sleep(time.Second * 10)
		}
	}()
	s := status{}

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Main Menu"
	update(s)
	//time.Sleep(time.Duration(rand.Int63()) % (time.Second * 15))

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Connecting to Open Match frontend"
	update(s)
	conn, err := rpc.GRPCClientFromConfig(cfg, "api.frontend")
	if err != nil {
		logger.Info(fmt.Sprintf("Encountered error: %s", err.Error()))
		panic(err)

	}
	defer conn.Close()
	fe := pb.NewFrontendClient(conn)
	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Creating Open Match Ticket"
	update(s)

	var ticketId string
	{
		req := &pb.CreateTicketRequest{
			Ticket: &pb.Ticket{
				Properties: structs.Struct{
					"name":      structs.String(name),
					"mode.demo": structs.Number(1),
					"player.level": structs.String(level),
					"player.exp": structs.Number(float64(exp)),
				}.S(),
			},
		}

		resp, err := fe.CreateTicket(ctx, req)
		if err != nil {
			logger.Info(fmt.Sprintf("Encountered error: %s", err.Error()))
			panic(err)
			
		}
		ticketId = resp.Ticket.Id
	}

	//logger.Info(fmt.Sprintf("Waiting match with ticket Id %s", ticketId))

	//////////////////////////////////////////////////////////////////////////////
	s.Status = fmt.Sprintf("Waiting match with ticket Id %s for level %s", ticketId, level)
	update(s)

	var assignment *pb.Assignment
	{
		req := &pb.GetAssignmentsRequest{
			TicketId: ticketId,
		}

		stream, err := fe.GetAssignments(ctx, req)
		for assignment.GetConnection() == "" {
			resp, err := stream.Recv()
			if err != nil {
				// For now we don't expect to get EOF, so that's still an error worthy of panic.
				panic(err)
			}

			assignment = resp.Assignment
		}

		err = stream.CloseSend()
		if err != nil {
			panic(err)
		}
	}

	//////////////////////////////////////////////////////////////////////////////
	s.Status = "Sleeping (pretend this is playing a match...)"
	s.Assignment = assignment
	update(s)

	time.Sleep(time.Second * 600)
}
