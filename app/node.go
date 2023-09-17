package app

import (
	"context"

	"github.com/bahner/go-space/config"
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	node.Node
}

var n node.Node

func StartApplication(ctx context.Context) {

	nodeName := config.NodeName
	nodeCookie := config.NodeCookie
	appName := config.AppName

	log.Infof("Starting %s Erlang Application node: %s (%s)\n", appName, nodeName, nodeCookie)

	var options node.Options
	var err error
	var process gen.Process

	// Create applications that must be started
	apps := []gen.ApplicationBehavior{
		createApplication(ctx),
	}
	options.Applications = apps

	// Starting node
	n, err = ergo.StartNodeWithContext(ctx, nodeName, nodeCookie, options)
	if err != nil {
		panic(err)
	}

	// Starting applications
	process, err = n.Spawn("space", gen.ProcessOptions{}, createSpace(ctx))
	if err != nil {
		panic(err)
	}

	log.Infof("Started process %q with PID %s.", process.Name(), process.Self())
}

func getNode() node.Node {
	return n
}
