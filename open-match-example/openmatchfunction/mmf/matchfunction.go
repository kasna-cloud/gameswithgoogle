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

// Package mmf provides a sample match function that uses the GRPC harness to set up
// the match making function as a service. This sample is a reference
// to demonstrate the usage of the GRPC harness and should only be used as
// a starting point for your match function. You will need to modify the
// matchmaking logic in this function based on your game's requirements.
package mmf

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	mmfHarness "open-match.dev/open-match/pkg/harness/function/golang"
	"open-match.dev/open-match/pkg/pb"
	
)

var (
	matchName = "10-player-based-match-with-red-blue-teams"
	logger          = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "mmf.pool",
	})
)

// MakeMatches is where your custom matchmaking logic lives.
// This is the core match making function that will be triggered by Open Match to generate matches.
// The goal of this function is to generate predictable matches that can be validated without flakyness.
// This match function loops through all the pools and generates one match per pool aggregating all players
// in that pool in the generated match.
func MakeMatches(params *mmfHarness.MatchFunctionParams) ([]*pb.Match, error) {
	var result []*pb.Match
	// logger.Info(fmt.Printf("%#v", params.PoolNameToTickets))
	for pool, tickets := range params.PoolNameToTickets {
		matchSize:=10
		for i := 0; i < len(tickets); i += matchSize {
			end := i + matchSize
			if end > len(tickets) {
				end = len(tickets)
			}
			logger.Info(fmt.Sprintf("totalcount for %s is %d tickets", pool, len(tickets)))
			matchtickets := tickets[i:end]
			if len(matchtickets) == 10 {
				logger.Info(fmt.Sprintf("count for %s is %d tickets", pool, len(matchtickets)))
				rosterblue := &pb.Roster{Name: "blue"}
				rosterred := &pb.Roster{Name: "red"}
				for _, ticket := range matchtickets[:5] {
					logger.Info(fmt.Sprintf("adding %s to blue", ticket.GetId()))
					rosterblue.TicketIds = append(rosterblue.GetTicketIds(), ticket.GetId())
				}
				for _, ticket := range matchtickets[5:] {
					logger.Info(fmt.Sprintf("adding %s to red", ticket.GetId()))
					rosterred.TicketIds = append(rosterred.GetTicketIds(), ticket.GetId())
				}
	
				result = append(result, &pb.Match{
					MatchId:       xid.New().String(),
					MatchProfile:  params.ProfileName,
					MatchFunction: matchName,
					Tickets:       matchtickets,
					Rosters:       []*pb.Roster{rosterblue,rosterblue},
				
				})		
			}
			
		}
	}
	return result, nil
}
