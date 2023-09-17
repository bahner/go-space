package keeper

import (
	"github.com/bahner/go-space/global"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/secrets"
	"gocloud.dev/secrets/hashivault"
)

// Create a secrets keeper. Each keeper is a wrapper around a single key in a
// HashiCorp Vault. The key is used to encrypt and decrypt secrets.

// The keeper is safe to use from multiple goroutines.

// The keeper must be closed when it is no longer needed.

// Each topic should have its own keeper. The key should be a unique value,
// which is the IPNS name.

// func New(client *api.Client, key string) *secrets.Keeper {
func New(key string) *secrets.Keeper {

	client := global.GetVaultClient()

	k := hashivault.OpenKeeper(client, key, nil)
	log.Infof("Created secrets keeper for key %s", key)

	return k
}
