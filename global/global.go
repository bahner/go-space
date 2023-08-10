package global

import (
	"context"
	"sync"
)

var (
	globalWg *sync.WaitGroup
	err      error
)

func StartServices(ctx context.Context) {

	// Init vault
	globalWg.Add(1)
	VaultClient, err = initVaultClient(ctx, vaultAddr, vaultToken)
	if err != nil {
		log.Fatal(err)
	}

	// Init pubsub
	globalWg.Add(1)
	PubSubService, err = startPubSubService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Waiting for global services to start")
	globalWg.Wait()
	log.Info("Global services started")

}
