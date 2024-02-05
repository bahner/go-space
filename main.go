package main

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-space/app"
	"github.com/spf13/pflag"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Init config and logger
	pflag.Parse()
	config.Init("space")
	config.InitLogging()
	config.InitP2P()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Fatalf("Error initialising P2P: %v", err)
	}
	log.Infof("My node ID is %s", p.Node.ID().String())

	p.DiscoverPeers()
	go p.DiscoveryLoop(context.Background())

	// Start application
	app.StartApplication(p)

	select {}
}
