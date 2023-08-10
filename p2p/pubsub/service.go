package pubsub

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/p2p/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Service struct {
	Sub  *pubsub.PubSub
	Host *host.P2pHost
}

func New(host *host.P2pHost) *Service {
	return &Service{
		Host: host,
	}
}

func (p *Service) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	log.Debug("Starting pubsub service...")
	p.Sub, err = pubsub.NewGossipSub(ctx, p.Host.Node)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p PubSub service started.")
}
