package app

import (
	"context"

	"github.com/bahner/go-space/config"
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	node.Node
}

var n node.Node

func StartApplication() {

	ctx := context.Background()

	nodeName := viper.GetString("node.name")
	nodeCookie := viper.GetString("node.cookie")
	appName := config.NAME

	log.Infof("Starting %s Erlang Application node: %s (%s)", appName, nodeName, nodeCookie)

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
		panic(err)
	}

	// Starting applications
	process, err = n.Spawn("space", gen.ProcessOptions{}, new(SPACE))
	if err != nil {
		panic(err)
	}

	log.Infof("Started process %q with PID %s.", process.Name(), process.Self())
}

func getNode() node.Node {
	return n
}
