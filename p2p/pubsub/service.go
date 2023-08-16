package pubsub

import (
	"context"

	"github.com/bahner/go-myspace/p2p/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Sub  *pubsub.PubSub
	Host *host.Host
}

func New(host *host.Host) *Service {
	return &Service{
		Host: host,
	}
}

func (p *Service) Start(ctx context.Context) {
	var err error

	log.Debug("Starting pubsub service...")
	p.Sub, err = pubsub.NewGossipSub(ctx, p.Host.Node)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p PubSub service started.")
}
