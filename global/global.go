package global

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/p2p/pubsub"
)

var (
	wg  *sync.WaitGroup
	err error
)

func InitGlobalResources(ctx context.Context) {

	VaultClient, err = initVaultClient(ctx, wg, vaultAddr, vaultToken)
	if err != nil {
		log.Fatal(err)
	}

	PubSubService = pubsub.New()
	PubSubService.Init(ctx)
}

func StartServices(ctx context.Context) {

	wg = &sync.WaitGroup{}

	// Init pubsub
	wg.Add(1)
	PubSubService, err = startPubSubService(ctx, wg)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Waiting for global services to start")
	wg.Wait()
	log.Info("Global services started")

}
