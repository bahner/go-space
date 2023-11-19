package ps

import (
	"context"
	"sync"

	"github.com/bahner/go-space/p2p/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var pubSubService *pubsub.PubSub

func InitPubSubService(ctx context.Context, wg *sync.WaitGroup, rendezvous string) {

	defer wg.Done()

	log.Info("Initializing global resources")

	// This needs no identity. It's anonymoous.
	// BUt maybe we should give it persistence later.
	host, err := host.New()
	if err != nil {
		log.Fatalf("Failed to create libp2p host: %v", err)
	}
	// vaultAddr := config.VaultAddr
	// vaultToken := config.VaultToken

	host.StartPeerDiscovery(ctx, rendezvous)

	pubSubService, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatalf("Failed to create pubsub: %v", err)
	}
	log.Info("Global resources initialized")

}

func GetService() *pubsub.PubSub {
	return pubSubService
}
