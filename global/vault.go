package global

import (
	"context"

	"github.com/hashicorp/vault/api"
	"gocloud.dev/secrets/hashivault"
)

func initVaultClient(ctx context.Context, addr string, token string) (*api.Client, error) {

	defer globalWg.Done()

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
