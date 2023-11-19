package global

import (
	"github.com/hashicorp/vault/api"
)

// func initVaultClient(ctx context.Context, wg *sync.WaitGroup, addr string, token string) error {
// 	defer wg.Done()

// 	client, err := hashivault.Dial(ctx, &hashivault.Config{
// 		Token: token,
// 		APIConfig: api.Config{
// 			Address: addr,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	vaultClient = client // Setting the package-level variable

// 	log.Info("Vault client initialized")

// 	return nil
// }

func GetVaultClient() *api.Client {
	return vaultClient
}
