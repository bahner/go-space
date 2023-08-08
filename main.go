package main

import (
	"context"
	"log"
	"myspace-pubsub/app"
	"myspace-pubsub/config"
	"myspace-pubsub/p2p"
)

func main() {

	ctx := context.Background()

	config.InitLogging()

	go p2p.StartPubSubService(ctx)

	n := app.StartApplication(ctx)

	status := n.IsAlive()
	log.Printf("Node is alive: %v\n", status)

	stats := n.Stats()
	log.Printf("Node stats: %v\n", stats)

	select {}
}
