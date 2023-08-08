package app

import (
	"context"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/node"
)

type Node struct {
	node.Node
}

func nodeStart(ctx context.Context) node.Node {

	log.Infof("Starting %s Erlang node: %s (%s)\n", appName, *nodeName, *nodeCookie)
	appNode, err := ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, node.Options{})
	if err != nil {
		panic(err)
	}

	log.Info("Application node started sucessfully.")

	return appNode
}

func nodeInit(ctx context.Context) error {

	n = nodeStart(ctx)

	return nil
}
