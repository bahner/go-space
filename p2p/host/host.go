package host

import (
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

type P2pHost struct {
	host.Host
}

func New(opts ...libp2p.Option) (*P2pHost, error) {

	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("p2p_host: Failed to create libp2p host: %v", err)
	}
	log.Info("libp2p node created: ", h.ID().String())

	return &P2pHost{h}, nil
}

func (h *P2pHost) StartPeerDiscovery(ctx context.Context, rendezvous string) {

	log.Debug("Starting peer discovery...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go discoverDHTPeers(ctx, wg, h, rendezvous)
	go discoverMDNSPeers(ctx, wg, h, rendezvous)
	wg.Wait()
}
