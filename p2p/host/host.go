package host

import (
	"context"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
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

func (h *P2pHost) StartPeerDiscovery(ctx context.Context, rendezvous string, serviceName string) {

	log.Debug("Starting peer discovery...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go DiscoverDHTPeers(ctx, wg, h.Node, rendezvous)
	go DiscoverMDNSPeers(ctx, wg, h.Node, serviceName)
	wg.Wait()
}
