package pubsub

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func StartPubSubService(ctx context.Context) {

	// Start libp2p node
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")
	go DiscoverPeers(ctx, h, *rendezvous)

	// Start pubsub service
	log.Debug("Starting pubsub service...")
	PubSubService, err = pubsub.NewGossipSub(ctx, h)
	if err != nil {
		// This is fatal because without pubsub, the app is useless.
		log.Fatal(err)
	}
}
