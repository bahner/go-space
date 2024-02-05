package main

import (
	"context"

	"github.com/bahner/go-space/app"
	"github.com/spf13/pflag"

	"github.com/bahner/go-space/config"
)

func main() {

	ctx := context.Background()
	// discoveryCtx, _ := context.WithTimeout(ctx, 60*time.Second)

	// Init config and logger
	pflag.Parse()
	config.Init()

	// n, ps, err := p2p.Init(discoveryCtx)

	// Start application
	app.StartApplication(ctx)

	select {}
}
