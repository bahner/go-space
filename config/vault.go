package config

import (
	"context"

	"github.com/hashicorp/vault/api"
	"gocloud.dev/secrets/hashivault"
)

// The reason for initializing Vault in config is that all
// external dependencies should be initialized in config.

// Although not safer, no need to keep passing around the
// tokens and addresses.

func initVaultClient(ctx context.Context, addr string, token string) (*api.Client, error) {
	client, err := hashivault.Dial(ctx, &hashivault.Config{
		Token: token,
		APIConfig: api.Config{
			Address: addr,
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
