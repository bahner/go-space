package config

import (
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	defaultNodeCookie = env.Get("MYSPACE_NODE_COOKIE", "myspace")
	defaultNodeName   = env.Get("MYSPACE_NODE_NAME", "pubsub@localhost")
	defaultRendezvous = env.Get("MYSPACE_LIBP2P_RENDEZVOUS", "myspace")
	defaultLogLevel   = env.Get("MYSPACE_LOG_LEVEL", "error")
)

var (
	Version     = "0.0.1"
	AppName     = "myspace-pubsub"
	Description = "The Myspace Pubsub Application"
	NodeCookie  = flag.String("nodecookie", defaultNodeCookie, "Secret shared by all erlang nodes in the cluster")
	NodeName    = flag.String("nodename", defaultNodeName, "Name of the erlang node")
	Log         = logrus.New()
	Rendezvous  = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	LogLevel    = flag.String("loglevel", defaultLogLevel, "Loglevel to use for application")
)

func InitLogging() {
	level, err := logrus.ParseLevel(*LogLevel)
	if err != nil {
		Log.Fatal(err)
	}
	Log.SetLevel(level)
}
