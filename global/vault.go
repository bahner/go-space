package global

import (
	"context"
	"sync"

	"github.com/hashicorp/vault/api"
	"gocloud.dev/secrets/hashivault"
)

func initVaultClient(ctx context.Context, wg *sync.WaitGroup, addr string, token string) (*api.Client, error) {

	defer wg.Done()

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
