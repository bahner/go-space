package main

import (
	"context"
	"myspace-pubsub/app"
	"myspace-pubsub/config"
	"myspace-pubsub/p2p"
)

var (
	ctx context.Context
	n   app.Node
)

func main() {

	ctx = context.Background()

	config.InitLogging()

	go p2p.StartPubSubService(ctx)

	// n, _ := app.Start(ctx)
	go app.Start(ctx)

	// go startMyspace(n)

	select {}
}
