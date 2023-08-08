package app

import (
	"context"

	"myspace-pubsub/config"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/node"
)

type Node struct {
	node.Node
}

var (
	log        = config.Log
	nodeCookie = config.NodeCookie
	nodeName   = config.NodeName
)

func Start(ctx context.Context) (node.Node, error) {

	log.Infof("Starting Erlang node: %s (%s)\n", *nodeName, *nodeCookie)
	n, err := ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, node.Options{})
	if err != nil {
		panic(err)
	}

	log.Info("done.")

	return n, nil
}
