package p2p

import (
	"context"
	"fmt"

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

func discoverMDNSPeers(ctx context.Context, h host.Host, rendezvous string) chan peer.AddrInfo {

	anyConnected := false
	for !anyConnected {

		log.Info("Starting MDNS peer discovery.")
		peerChan := initMDNS(h, rendezvous)

		for { // allows multiple peers to join
			peer := <-peerChan // will block until we discover a peer
			fmt.Println("Found MDNS peer:", peer.ID.Pretty(), ", connecting")

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				log.Infof("Connected to DHT peer: %s", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}

	log.Info("MDNS peer discovery successful.")

	return nil
}
