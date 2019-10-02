package main

import (
	"open-match-example/openmatchclient/demo"
	"open-match-example/openmatchclient/demo/components"
	"open-match-example/openmatchclient/demo/components/clients"
	"open-match-example/openmatchclient/demo/components/director"
	"open-match-example/openmatchclient/demo/components/uptime"
)

func main() {
	demo.Run(map[string]func(*components.DemoShared){
		"uptime":   uptime.Run,
		"clients":  clients.Run,
		"director": director.Run,
	})
}
