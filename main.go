package main

import (
	"context"

	"github.com/bahner/go-space/app"

	"github.com/bahner/go-space/config"
)

func main() {

	ctx := context.Background()
	// discoveryCtx, _ := context.WithTimeout(ctx, 60*time.Second)

	// Init config and logger
	config.Init(ctx)

	// n, ps, err := p2p.Init(discoveryCtx)

	// Start application
	app.StartApplication(ctx)

	select {}
}
