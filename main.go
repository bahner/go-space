package main

import (
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

	noe, err := p2p.Init(nil)
	if err != nil {
		log.Fatalf("Error initialising P2P: %v", err)
	}

	// Start application
	app.StartApplication(noe)

	select {}
}
