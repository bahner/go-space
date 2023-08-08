package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// discoverPeers performs peer discovery using the DHT and connects to discovered peers
func DiscoverPeers(ctx context.Context, h host.Host, rendezvousString string) error {
	dhtInstance, err := initDHT(ctx, h)
	if err != nil {
		return err
	}

	// Set up an mDNS service on the libp2p host
	// ser, err := mdns.NewMdnsService(ctx, h, time.Second*5, rendezvousString)
	// ser := mdns.NewMdnsService(ctx, h, time.Second*5, rendezvousString)
	// ser.RegisterNotifee(&discovery.Notifee{})
	// if err != nil {
	// 	panic(err)
	// }

	// The service will run in the background printing discovered peers to the console
	// ser.RegisterNotifee(&discovery.Notifee{})

	routingDiscovery := drouting.NewRoutingDiscovery(dhtInstance)
	dutil.Advertise(ctx, routingDiscovery, rendezvousString)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		log.Info("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		if err != nil {
			return fmt.Errorf("peer discovery error: %w", err)
		}

		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // Skip self connection
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Errorf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				fmt.Println("Connected to:", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}

	log.Info("Peer discovery complete")

	return nil
}
