package config

import (
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	defaultNodeCookie      = env.Get("MYSPACE_NODE_COOKIE", "myspace")
	defaultNodeName        = env.Get("MYSPACE_NODE_NAME", "pubsub@localhost")
	defaultRendezvous      = env.Get("MYSPACE_LIBP2P_RENDEZVOUS", "myspace")
	defaultLogLevel        = env.Get("MYSPACE_LOG_LEVEL", "error")
	defaultMyspaceNodeName = env.Get("MYSPACE_MYSPACE_NODE_NAME", "myspace@localhost")
	defaultVaultAddr       = env.Get("MYSPACE_VAULT_ADDR", "http://localhost:8200")
)

var (
	Version     = "0.0.1"
	AppName     = "myspace-pubsub"
	Description = "The Myspace Pubsub Application"

	Log = logrus.New()

	LogLevel        = flag.String("loglevel", defaultLogLevel, "Loglevel to use for application")
	MyspaceNodeName = flag.String("myspace_nodename", defaultMyspaceNodeName, "Name of the node running the actual Myspace")
	NodeCookie      = flag.String("nodecookie", defaultNodeCookie, "Secret shared by all erlang nodes in the cluster")
	NodeName        = flag.String("nodename", defaultNodeName, "Name of the erlang node")
	Rendezvous      = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	VaultAddr       = flag.String("vaultaddr", defaultVaultAddr, "Address of the vault server")
	VaultToken      = flag.String("vaulttoken", "", "Token to use to authenticate with the vault server. This is required.")
)

func InitLogging() {
	level, err := logrus.ParseLevel(*LogLevel)
	if err != nil {
		Log.Fatal(err)
	}
	Log.SetLevel(level)
}
