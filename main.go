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

	// Init config and common services
	config.Init(ctx)

	// Start p2p node and services
	go p2p.StartPubSubService(ctx)

	// Start Erlang node and application
	n := app.StartApplication(ctx)

	status := n.IsAlive()
	log.Printf("Node is alive: %v\n", status)

	stats := n.Stats()
	log.Printf("Node stats: %v\n", stats)

	select {}
}
