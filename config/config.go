package config

import (
	"context"
	"flag"

	"github.com/bahner/go-myspace/p2p/key"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	// Erlang application config
	Version     = "0.0.1"
	AppName     = "go-myspace"
	Description = "Myspace node written in go to handle libp2p functionality."

	// Package internal config
	VaultAddr  string = env.Get("GO_MYSPACE_VAULT_ADDR", "http://localhost:8200")
	VaultToken string = env.Get("GO_MYSPACE_VAULT_TOKEN", "myspace")

	// Global config
	LogLevel        string = env.Get("GO_MYSPACE_LOG_LEVEL", "error")
	MyspaceNodeName string = env.Get("GO_MYSPACE_MYSPACE_NODE_NAME", "myspace@localhost")
	NodeCookie      string = env.Get("GO_MYSPACE_NODE_COOKIE", "myspace")
	NodeName        string = env.Get("GO_MYSPACE_NODE_NAME", "pubsub@localhost")
	Rendezvous      string = env.Get("GO_MYSPACE_RENDEZVOUS", "myspace")
	Secret          string = env.Get("GO_MYSPACE_P2P_IDENTITY", "")
	ServiceName     string = env.Get("GO_MYSPACE_SERVICE_NAME", "myspace")
)

func Init(ctx context.Context) {

	// Flags - user configurations
	flag.StringVar(&LogLevel, "loglevel", LogLevel, "Loglevel to use for application")
	flag.StringVar(&MyspaceNodeName, "myspace_nodename", MyspaceNodeName, "Name of the node running the actual Myspace")
	flag.StringVar(&NodeCookie, "nodecookie", NodeCookie, "Secret shared by all erlang nodes in the cluster")
	flag.StringVar(&NodeName, "nodename", NodeName, "Name of the erlang node")
	flag.StringVar(&Rendezvous, "rendezvous", Rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&Secret, "identity", Secret, "Base58 encoded secret key used to identify libp2p node for persistency.")
	flag.StringVar(&ServiceName, "servicename", ServiceName, "serviceName to use for MDNS discovery")
	flag.StringVar(&VaultAddr, "vaultaddr", VaultAddr, "Address of the vault server")
	flag.StringVar(&VaultToken, "vaulttoken", VaultToken, "Token to use to authenticate with the vault server. This is required.")

	generate := flag.Bool("generate", false, "Generate a new identity")

	flag.Parse()

	if *generate {
		key.PrintEncodedKeyAndExit()
	}

	// Init logger
	level, err := log.ParseLevel(LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

}
