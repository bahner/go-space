package global

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
	"github.com/hashicorp/vault/api"
)

var (
	err           error
	vaultClient   *api.Client
	pubSubService *pubsub.Service
)

func StartServices(ctx context.Context) {

	log := config.GetLogger()
	vaultAddr := config.VaultAddr
	vaultToken := config.VaultToken

	log.Info("Initializing global resources")

	wg := &sync.WaitGroup{}

	p2phost := host.New()
	p2phost.Init(ctx)
	wg.Add(1)
	p2phost.StartPeerDiscovery(ctx, wg)
	log.Info("Waiting for P2P host to start")
	wg.Wait()
	log.Info("P2P host started")

	pubSubService = pubsub.New(p2phost)
	wg.Add(1)
	pubSubService.Start(ctx, wg)
	log.Info("Waiting for pubsub service to start")
	wg.Wait()
	log.Info("Pubsub service started")

	// Initialize vault client
	vaultClient, err = initVaultClient(ctx, vaultAddr, vaultToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Vault client initialized")

}

func GetVaultClient() *api.Client {
	return vaultClient
}

func GetPubSubService() *pubsub.Service {
	return pubSubService
}
