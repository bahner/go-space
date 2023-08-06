package main

import (
	"context"
	"flag"
	"log"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	"go.deanishe.net/env"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var (
	defaultRendezvous = env.Get("MYSPACE_LIBP2P_RENDEZVOUS", "myspace")
	defaultLogLevel   = env.Get("MYSPACE_LIBP2P_LOG_LEVEL", "error")
	defaultNodeCookie = env.Get("MYSPACE_NODE_COOKIE", "myspace")
	defaultNodeName   = env.Get("MYSPACE_NODE_NAME", "go@localhost")
)

var (
	// h          host.Host
	ps         *pubsub.PubSub
	logger     = logging.Logger("myspace")
	logLevel   = flag.String("loglevel", defaultLogLevel, "Log level for libp2p")
	rendezvous = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	nodeCookie = flag.String("nodecookie", defaultNodeCookie, "Secret shared by all erlang nodes in the cluster")
	nodeName   = flag.String("nodename", defaultNodeName, "Name of the erlang node")
	n          node.Node
	ctx        context.Context
)

func main() {

	ctx = context.Background()

	// libp2p node
	logging.SetLogLevel("myspace", *logLevel)
	logger.Info("Starting myspace libp2p pubsub server...")

	// Start libp2p h
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	go DiscoverPeers(ctx, h, *rendezvous)

	// Start pubsub service
	ps, err = pubsub.NewGossipSub(ctx, h)
	if err != nil {
		// This is fatal because without pubsub, the app is useless.
		log.Fatal(err)
	}

	// Erlang node
	n, err = ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, node.Options{})
	// This is fatal because without an erlang node, the app is useless.
	if err != nil {
		log.Fatal(err)
	}

	spawnAndRegisterRoom("lobby")

	select {}
}

func spawnAndRegisterRoom(roomID string) {
	process, err := n.Spawn(roomID, gen.ProcessOptions{}, createRoom(roomID), roomID)
	if err != nil {
		log.Fatal(err)
	}
	n.RegisterName(roomID, process.Self())
}
