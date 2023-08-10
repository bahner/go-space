package host

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
func discoverMDNSPeers(ctx context.Context, wg *sync.WaitGroup, h host.Host, rendezvous string) chan peer.AddrInfo {

	defer wg.Done()

	log := config.GetLogger()

	anyConnected := false
	for !anyConnected {
		log.Info("Starting MDNS peer discovery.")
		peerChan := initMDNS(h, rendezvous)

		// Keep the loop running until you've connected to a peer
		for !anyConnected {
			peer := <-peerChan // will block until we discover a peer
			log.Infof("Found MDNS peer: %s connecting", peer.ID.Pretty())

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				log.Infof("Connected to MDNS peer: %s", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}

	log.Info("MDNS peer discovery successful.")
	return nil
}
