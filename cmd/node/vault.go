package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-ma/key/set"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	defaultVaultAddr = "http://localhost:8200"
	devVaultToken    = "space"
	keysetPath       = "keysets/"
)

var (
	once     sync.Once
	vaultApi *api.Client
)

func init() {

	pflag.String("vault-addr", defaultVaultAddr, "Vault address")
	viper.BindPFlag("vault.addr", pflag.Lookup("vault-addr"))
	viper.SetDefault("vault.addr", defaultVaultAddr)

	pflag.String("vault-token", "", "Vault token")
	viper.BindPFlag("vault.token", pflag.Lookup("vault-token"))
	viper.SetDefault("vault.token", devVaultToken)
}

func VaultClient() (*api.Client, error) {

	once.Do(func() {

		config := api.DefaultConfig()
		config.Address = viper.GetString("vault.addr")

		client, err := api.NewClient(config)
		if err != nil {
			log.Errorf("Error creating vault client: %s", err)
		}

		client.SetToken(viper.GetString("vault.token"))

		vaultApi = client

	})

	return vaultApi, nil
}

// Store the keyset the id is the fragment, the nick, not the full did
func storeKeyset(keyset set.Keyset) error {

	log.Debugf("Storing keyset: %s", keyset.DID.Fragment)
	id := keyset.DID.Fragment

	keysetString, err := keyset.Pack()
	if err != nil {
		return fmt.Errorf("error packing keyset: %s", err)
	}

	log.Debugf("Packaged keyset: %s", keysetString)
	return storeKeysetString(id, keysetString)
}

func storeKeysetString(id string, keyset string) error {

	log.Debugf("Storing keyset string in vault: %s", id)
	client, err := VaultClient()
	if err != nil {
		return err
	}

	secret := map[string]interface{}{
		"multibaseKeyset": keyset,
	}

	log.Debugf("Storing secret: %s", secret)
	_, err = client.KVv2("secret").Put(context.Background(), keysetPath+id, secret)
	if err != nil {
		return fmt.Errorf("error storing keyset: %s", err)
	}

	log.Infof("Stored keyset: %s", id)
	return nil
}
func retrieveKeyset(id string) (set.Keyset, error) {

	keysetString, err := retrieveKeysetString(id)
	if err != nil {
		return set.Keyset{}, err
	}

	return set.Unpack(keysetString)
}

func retrieveKeysetString(id string) (string, error) {

	client, err := VaultClient()
	if err != nil {
		return "", err
	}

	secret, err := client.KVv2("secret").Get(context.Background(), keysetPath+id)
	if err != nil {
		return "", err
	}

	keyset, ok := secret.Data["multibaseKeyset"].(string)
	if !ok {
		return "", err
	}

	return keyset, nil
}

func getOrCreateKeysetFromVault(id string) (set.Keyset, error) {

	keyset, err := retrieveKeyset(id)
	if err != nil {
		if err != api.ErrSecretNotFound {
			return set.Keyset{}, fmt.Errorf("error retrieving keyset: %s", err)
		}

		keyset, err = set.GetOrCreate(id)
		if err != nil {
			return set.Keyset{}, fmt.Errorf("error getting or creating keyset: %s", err)
		}

		err = storeKeyset(keyset)
		if err != nil {
			return set.Keyset{}, fmt.Errorf("error storing keyset: %s", err)
		}
	}

	return keyset, nil
}
