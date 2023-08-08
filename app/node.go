package node

import (
	"context"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/node"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	node.Node
}

func Start(ctx context.Context) (node.Node, error) {

	log.Infof("Starting Erlang node: %s (%s)\n", *nodeName, *nodeCookie)
	n, err := ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, node.Options{})
	if err != nil {
		panic(err)
	}

	log.Info("done.")

	return n, nil
}
