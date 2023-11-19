package global

import (
	"context"
	"sync"

	"github.com/bahner/go-space/config"
	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-space/p2p/pubsub"
)

func initPubSubService(ctx context.Context, wg *sync.WaitGroup, host *host.P2pHost) {

	// Start libp2p node and discover peers
	host.StartPeerDiscovery(ctx, wg, config.Rendezvous)

	pubSubService = pubsub.New(ctx, host)

}

func GetPubSubService() *pubsub.Service {
	return pubSubService
}
