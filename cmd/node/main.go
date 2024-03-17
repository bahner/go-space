package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	// Init config and logger
	initConfig()

	p, err := initP2P()
	if err != nil {
		log.Fatalf("Error initialising P2P: %v", err)
	}

	// Init of actor requires P2P to be initialized
	a := initActorOrPanic()

	go p.StartDiscoveryLoop(ctx)
	go a.Subscribe(ctx, a.Entity)

	// Start application
	StartApplication(p)

	select {}
}
