package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

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
	ps         *pubsub.PubSub
	libp2pLog  = logging.Logger("myspace")
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
	libp2pLog.Info("Starting myspace libp2p pubsub server...")

	// Start libp2p h
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")
	go DiscoverPeers(ctx, h, *rendezvous)

	// Start pubsub service
	log.Debug("Starting pubsub service...")
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

	log.Info("Starting myspace dispatcher")
	n.Spawn("myspace", gen.ProcessOptions{}, createMyspace(), "myspace")

	select {}
}
