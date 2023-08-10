package pubsub

import (
	"context"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type Service struct {
	Node   host.Host
	PubSub *pubsub.PubSub
}

func New() *Service {
	return &Service{}
}

func (p *Service) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	log.Info("Starting libp2p node...")
	p.Node, err = libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", p.Node.ID().Pretty(), " ", p.Node.Addrs())

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")
	var discoveryWg sync.WaitGroup
	discoveryWg.Add(2)

	go discoverDHTPeers(ctx, &discoveryWg, p.Node, rendezvous)
	go discoverMDNSPeers(ctx, &discoveryWg, p.Node, rendezvous)

	// Wait for both discovery processes to complete
	discoveryWg.Wait()

	log.Debug("Starting pubsub service...")
	p.PubSub, err = pubsub.NewGossipSub(ctx, p.Node)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p PubSub service started.")
}
