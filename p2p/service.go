package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var (
	PubSubService *pubsub.PubSub
)

func StartPubSubService(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Info("Starting libp2p node...")
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.ID().Pretty(), " ", h.Addrs())

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")
	var discoveryWg sync.WaitGroup
	discoveryWg.Add(2)

	go discoverDHTPeers(ctx, &discoveryWg, h, rendezvous)
	go discoverMDNSPeers(ctx, &discoveryWg, h, rendezvous)

	// Wait for both discovery processes to complete
	discoveryWg.Wait()

	log.Debug("Starting pubsub service...")
	PubSubService, err = pubsub.NewGossipSub(ctx, h)
	if err != nil {
		// This is fatal because without pubsub, the app is useless.
		log.Fatal(err)
	}
	log.Info("PubSub service started.")
}
