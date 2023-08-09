package p2p

import (
	"context"
	"fmt"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func initDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT, nil
}

// discoverPeers performs peer discovery using the DHT and connects to discovered peers
func discoverDHTPeers(ctx context.Context, h host.Host, rendezvousString string) error {
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
		log.Info("Starting DHT peer discovery.")
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
				log.Debugf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				log.Infof("Connected to: %s", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}

	log.Info("Peer discovery complete")

	return nil
}
