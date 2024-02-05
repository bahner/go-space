package main

import (
	"github.com/bahner/go-space/app"
	"github.com/spf13/pflag"

	"github.com/bahner/go-ma-actor/config"
)

func main() {

	// Init config and logger
	pflag.Parse()
	config.Init("space")

	// n, ps, err := p2p.Init(discoveryCtx)

	// Start application
	app.StartApplication()

	select {}
}
