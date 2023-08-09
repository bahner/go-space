package keeper

import (
	"context"

	"gocloud.dev/secrets"
)

func Encrypt(keeper *secrets.Keeper, data []byte) ([]byte, error) {

	ctx := context.Background()

	result, err := keeper.Encrypt(ctx, data)
	if err != nil {
		return nil, err
	}

	defer keeper.Close()

	return result, nil
}

func Decrypt(keeper *secrets.Keeper, data []byte) ([]byte, error) {

	ctx := context.Background()

	result, err := keeper.Decrypt(ctx, data)
	if err != nil {
		return nil, err
	}

	defer keeper.Close()

	return result, nil
}
