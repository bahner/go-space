package main

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p"
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	n node.Node
	p *p2p.P2P
)

func StartApplication(_p2p *p2p.P2P) {

	ctx := context.Background()
	p = _p2p

	nodeName := viper.GetString("node.name")
	nodeCookie := viper.GetString("node.cookie")

	log.Infof("Starting %s Erlang Application node: %s (%s)", NAME, nodeName, nodeCookie)

	var options node.Options
	var err error
	var process gen.Process

	// Create applications that must be started
	apps := []gen.ApplicationBehavior{
		new(Application),
	}
	options.Applications = apps

	// Starting node
	n, err = ergo.StartNodeWithContext(ctx, nodeName, nodeCookie, options)
	if err != nil {
		log.Fatalf("Failed to start node: %s", err)
	}

	// Starting applications
	process, err = n.Spawn("space", gen.ProcessOptions{}, new(SPACE))
	if err != nil {
		log.Fatalf("Failed to spawn space application: %s", err)
	}

	log.Infof("Started process %q with PID %s.", process.Name(), process.Self())
}

func getNode() node.Node {
	return n
}
