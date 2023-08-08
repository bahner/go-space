package main

import (
	"context"
	"myspace-pubsub/app"
	"myspace-pubsub/config"
	"myspace-pubsub/p2p"
)

func main() {

	ctx := context.Background()

	config.InitLogging()

	p2p.StartPubSubService(ctx)

	app.StartApplication(ctx)

	select {}
}
