package pubsub

import (
	"flag"

	logging "github.com/ipfs/go-log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.deanishe.net/env"
)

var (
	defaultRendezvous = env.Get("MYSPACE_LIBP2P_RENDEZVOUS", "myspace")
	defaultLogLevel   = env.Get("MYSPACE_LIBP2P_LOG_LEVEL", "error")
)

var (
	PubSubService *pubsub.PubSub
	libp2pLog     = logging.Logger("myspace")
	logLevel      = flag.String("loglevel", defaultLogLevel, "Log level for libp2p")
	rendezvous    = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
)
