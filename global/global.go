package global

import (
	"context"
	"sync"

	"github.com/bahner/go-space/config"
	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-sub-etha/libp2p/key"
	"github.com/bahner/go-space/p2p/pubsub"
	"github.com/hashicorp/vault/api"
	"github.com/libp2p/go-libp2p"
	log "github.com/sirupsen/logrus"
)

var (
	vaultClient   *api.Client
	pubSubService *pubsub.Service
)

func InitAndStartServices(ctx context.Context) {

	// Configure host
	s := config.Secret
	h := host.New()
	if config.Secret != "" {
		id, err := key.Unpack(s)
		if err != nil {
			log.Fatal(err)
		}
		h.AddOption(libp2p.Identity(id.PrivKey))
	}

	// Configure vault
	vaultAddr := config.VaultAddr
	vaultToken := config.VaultToken

	wg := &sync.WaitGroup{}
	wg.Add(2)

	log.Info("Initializing global resources")

	initPubSubService(ctx, wg, h)
	initVaultClient(ctx, wg, vaultAddr, vaultToken)

	log.Info("Waiting for global resources to be initialized ...")

	wg.Wait()

	log.Info("Global resources initialized")

}
