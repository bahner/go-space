package global

import (
	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/p2p/pubsub"
	vault "github.com/hashicorp/vault/api"
)

var (
	VaultClient   *vault.Client
	PubSubService *pubsub.Service
)

var (
	vaultAddr  = config.VaultAddr
	vaultToken = config.VaultToken
	log        = config.Log
)
