package host

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/config"
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

	log := config.GetLogger()

	var err error
	log.Info("Starting libp2p node...")
	h.Node, err = libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.Node.ID().Pretty())
}

func (h *P2pHost) StartPeerDiscovery(ctx context.Context, wg *sync.WaitGroup) {

	log := config.GetLogger()
	rendezvous := config.Rendezvous

	log.Debug("Starting peer discovery...")

	wg.Add(2)
	go discoverDHTPeers(ctx, wg, h.Node, rendezvous)
	go discoverMDNSPeers(ctx, wg, h.Node, rendezvous)
}
