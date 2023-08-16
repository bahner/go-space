package global

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
)

func initPubSubService(ctx context.Context, wg *sync.WaitGroup, host *host.Host) {

	defer wg.Done()

	// Start libp2p node and discover peers
	host.Init(ctx)
	host.StartPeerDiscovery(ctx, config.Rendezvous, config.ServiceName)

	pubSubService = pubsub.New(host)
	pubSubService.Start(ctx)

}

func GetPubSubService() *pubsub.Service {
	return pubSubService
}
