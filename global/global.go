package global

import (
	"context"
	"sync"

	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-space/p2p/pubsub"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

var (
	vaultClient   *api.Client
	pubSubService *pubsub.Service
)

func InitAndStartServices(ctx context.Context) {

	// This needs no identity. It's anonymoous.
	// BUt maybe we should give it persistence later.
	host, err := host.New()
	if err != nil {
		log.Fatalf("Failed to create libp2p host: %v", err)
	}
	// vaultAddr := config.VaultAddr
	// vaultToken := config.VaultToken

	wg := &sync.WaitGroup{}
	wg.Add(2)

	log.Info("Initializing global resources")

	initPubSubService(ctx, wg, host)
	// initVaultClient(ctx, wg, vaultAddr, vaultToken)

	log.Info("Waiting for global resources to be initialized ...")

	wg.Wait()

	log.Info("Global resources initialized")

}
