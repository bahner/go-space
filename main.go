package main

import (
	"context"
	"myspace/node"
	"myspace/pubsub"
)

var (
	ctx context.Context
	n   node.Node
)

func main() {

	ctx = context.Background()

	go pubsub.StartPubSubService(ctx)

	n, _ := node.Start(ctx)

	go startMyspace(n)

	select {}
}
