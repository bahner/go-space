package app

import (
	"context"
	"fmt"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
)

type Node struct {
	node.Node
}

func StartApplication(ctx context.Context) node.Node {

	log.Infof("Starting %s Erlang Application node: %s (%s)\n", appName, *nodeName, *nodeCookie)

	var options node.Options
	var err error
	var process gen.Process

	// Create applications that must be started
	apps := []gen.ApplicationBehavior{
		createApplication(ctx),
	}
	options.Applications = apps

	// Starting node
	n, err = ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, options)
	if err != nil {
		panic(err)
	}

	log.Info("Application node started sucessfully.")

	// Starting applications
	process, err = n.Spawn("myspace", gen.ProcessOptions{}, createMyspace(ctx))
	if err != nil {
		panic(err)
	}
	fmt.Printf("  process %q with PID %s is started\n", process.Name(), process.Self())

	return n
}
