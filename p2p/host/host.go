package host

import (
	"context"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2pHost struct {
	Node host.Host
}

func New() *P2pHost {
	return &P2pHost{}
}

func (h *P2pHost) Init(ctx context.Context) {
	var err error
	log.Info("Starting libp2p node...")
	h.Node, err = libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.Node.ID().Pretty())
}

func (h *P2pHost) StartPeerDiscovery(ctx context.Context) {
	log.Debug("Starting peer discovery...")
	var discoveryWg sync.WaitGroup
	discoveryWg.Add(2)

	go discoverDHTPeers(ctx, &discoveryWg, h.Node, rendezvous)
	go discoverMDNSPeers(ctx, &discoveryWg, h.Node, rendezvous)

	// Wait for both discovery processes to complete
	discoveryWg.Wait()
}
