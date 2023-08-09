package main

import (
	"context"
	"log"

	"github.com/bahner/go-myspace/p2p"

	"github.com/bahner/go-myspace/config"

	"github.com/bahner/go-myspace/app"
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
