package main

import (
	"github.com/bahner/go-space/app"
	"github.com/spf13/pflag"

	"github.com/bahner/go-space/config"
)

func main() {

	// Init config and logger
	pflag.Parse()
	config.Init()

	// n, ps, err := p2p.Init(discoveryCtx)

	// Start application
	app.StartApplication()

	select {}
}
