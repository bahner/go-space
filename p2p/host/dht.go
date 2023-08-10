package host

import (
	"context"
	"fmt"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

func initDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {

	log.Info("Initializing DHT.")

	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		log.Error("Failed to create Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT created.")
	}

	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		log.Error("Failed to bootstrap Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT bootstrap setup.")
	}

	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Warnf("Failed to convert bootstrap peer address: %v", err)
			continue
		}

		log.Debugf("Bootstrapping to peer: %s", peerinfo.ID.Pretty())

		go func(pInfo peer.AddrInfo) {

			log.Debugf("Attempting connection to peer: %s", pInfo.ID.Pretty())

			if err := h.Connect(ctx, pInfo); err != nil {
				log.Warnf("Bootstrap warning: %v", err)
			}
		}(*peerinfo)
	}

	log.Info("Kademlia DHT bootstrapped successfully.")
	return kademliaDHT, nil
}

func discoverDHTPeers(ctx context.Context, wg *sync.WaitGroup, h host.Host, rendezvousString string) error {

	defer wg.Done()

	log.Debug("Starting DHT route discovery.")

	dhtInstance, err := initDHT(ctx, h)
	if err != nil {
		return err
	}

	routingDiscovery := drouting.NewRoutingDiscovery(dhtInstance)
	dutil.Advertise(ctx, routingDiscovery, rendezvousString)

	log.Infof("Starting DHT peer discovery for rendezvous string: %s", rendezvousString)

	retryCount := 0

	for {

		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		if err != nil {
			return fmt.Errorf("peer discovery error: %w", err)
		}

		anyConnected := false
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // Skip self connection
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				log.Infof("Connected to DHT peer: %s", peer.ID.Pretty())
				anyConnected = true
			}
		}

		if anyConnected {
			break
		}
		retryCount++
		log.Debugf("Attempts #%d for peer discovery with rendezvous string: %s failed.", retryCount, rendezvousString)
	}

	log.Info("Peer discovery complete")
	return nil
}
