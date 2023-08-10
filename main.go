package main

import (
	"context"

	"github.com/bahner/go-myspace/global"
	"github.com/bahner/go-myspace/keeper"

	"github.com/bahner/go-myspace/config"
)

func main() {

	ctx := context.Background()

	// Init config and common services
	config.Init(ctx)

	log := config.Log

	// Start background services
	global.StartServices(ctx)

	_secret := []byte("secret")

	k := keeper.New("myspace")
	defer k.Close()

	safe_secret, err := keeper.Encrypt(k, _secret)
	if err != nil {
		log.Fatal(err)
	}
	log.Warnf("safe_secret: %v", safe_secret)

	select {}
}
