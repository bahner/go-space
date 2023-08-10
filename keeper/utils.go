package keeper

import (
	"context"

	"gocloud.dev/secrets"
)

func Encrypt(k *secrets.Keeper, data []byte) ([]byte, error) {

	ctx := context.Background()

	result, err := k.Encrypt(ctx, data)
	if err != nil {
		return nil, err
	}

	defer k.Close()

	return result, nil
}

func Decrypt(k *secrets.Keeper, data []byte) ([]byte, error) {

	ctx := context.Background()

	result, err := k.Decrypt(ctx, data)
	if err != nil {
		return nil, err
	}

	defer k.Close()

	return result, nil
}
