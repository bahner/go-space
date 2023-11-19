package main

import (
	"context"
	"sync"

	"github.com/bahner/go-space/app"
	"github.com/bahner/go-space/ps"

	"github.com/bahner/go-space/config"
)

func main() {

	ctx := context.Background()

	// Init config and logger
	config.Init(ctx)

	wg := &sync.WaitGroup{}
	// Start background services
	wg.Add(1)
	ps.InitPubSubService(ctx, wg, config.Rendezvous)
	wg.Wait()

	// Start application
	app.StartApplication(ctx)

	select {}
}
