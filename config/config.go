package config

import (
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	Version     = "0.0.1"
	AppName     = "go-myspace"
	Description = "Myspace node written in go to handle libp2p functionality."

	Log = logrus.New()

	LogLevel        string = env.Get("GO_MYSPACE_LOG_LEVEL", "error")
	MyspaceNodeName string = env.Get("GO_MYSPACE_MYSPACE_NODE_NAME", "myspace@localhost")
	NodeCookie      string = env.Get("GO_MYSPACE_NODE_COOKIE", "myspace")
	NodeName        string = env.Get("GO_MYSPACE_NODE_NAME", "pubsub@localhost")
	Rendezvous      string = env.Get("GO_MYSPACE_RENDEZVOUS", "myspace")
	VaultAddr       string = env.Get("GO_MYSPACE_VAULT_ADDR", "http://localhost:8200")
	VaultToken      string = env.Get("GO_MYSPACE_VAULT_TOKEN", "myspace")
)

func InitLogging() {

	flag.StringVar(&LogLevel, "loglevel", LogLevel, "Loglevel to use for application")
	flag.StringVar(&MyspaceNodeName, "myspace_nodename", MyspaceNodeName, "Name of the node running the actual Myspace")
	flag.StringVar(&NodeCookie, "nodecookie", NodeCookie, "Secret shared by all erlang nodes in the cluster")
	flag.StringVar(&NodeName, "nodename", NodeName, "Name of the erlang node")
	flag.StringVar(&Rendezvous, "rendezvous", Rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&VaultAddr, "vaultaddr", VaultAddr, "Address of the vault server")
	flag.StringVar(&VaultToken, "vaulttoken", VaultToken, "Token to use to authenticate with the vault server. This is required.")

	flag.Parse()

	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		Log.Fatal(err)
	}
	Log.SetLevel(level)
}
