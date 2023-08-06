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
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	defaultRendezvous = env.Get("MYSPACE_LIBP2P_RENDEZVOUS", "myspace")
	defaultLogLevel   = env.Get("MYSPACE_LIBP2P_LOG_LEVEL", "error")
	defaultNodeCookie = env.Get("MYSPACE_NODE_COOKIE", "myspace")
	defaultNodeName   = env.Get("MYSPACE_NODE_NAME", "go")
)

var (
	h          host.Host
	ps         *pubsub.PubSub
	logger     = logging.Logger("myspace")
	logLevel   = flag.String("loglevel", defaultLogLevel, "Log level for libp2p")
	rendezvous = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	nodeCookie = flag.String("nodecookie", defaultNodeCookie, "Secret shared by all erlang nodes in the cluster")
	nodeName   = flag.String("nodename", defaultNodeName, "Name of the erlang node")
)

func main() {

	ctx := context.Background()

	// libp2p node
	logging.SetLogLevel("myspace", *logLevel)
	logger.Info("Starting myspace libp2p pubsub server...")

	// Start libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	go DiscoverPeers(ctx, host, *rendezvous)

	// Start pubsub service
	ps, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		// This is fatal because without pubsub, the app is useless.
		log.Fatal(err)
	}

	// Erlang node
	node, err := ergo.StartNodeWithContext(ctx, *nodeName, *nodeCookie, node.Options{})
	// This is fatal because without an erlang node, the app is useless.
	if err != nil {
		log.Fatal(err)
	}

	supOpts := gen.Supervisor.Options{
		Strategy:  gen.SupervisorStrategyOneForOne,
		Intensity: 1,
		Period:    5,
	}

	// Supervisor
	supSpec := ergo.SupervisorSpec{
		Name: "gameSupervisor",
		Children: []ergo.SupervisorChildSpec{
			{
				Name: "goRoom",
				ChildGenServer: ergo.SupervisorChildGenServer{
					Args: []interface{}{"room1"},
					Func: func() ergo.GenServer { return &goRoom{} },
				},
			},
			{
				Name: "goAvatar",
				ChildGenServer: ergo.SupervisorChildGenServer{
					Args: []interface{}{"avatar1"},
					Func: func() ergo.GenServer { return &goAvatar{} },
				},
			},
		},
	}

	_, _ = node.Supervisor(supOpts, supSpec)
	select {}
}
