package pubsub

import (
	"context"

	"github.com/bahner/go-space/p2p/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Sub  *pubsub.PubSub
	Host *host.P2pHost
}

func New(ctx context.Context, h *host.P2pHost) *Service {

	log.Debug("Starting pubsub service...")
	sub, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p PubSub service started.")
	return &Service{
		Sub:  sub,
		Host: h,
	}
}
