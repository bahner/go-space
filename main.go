package main

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/keeper"
	"github.com/bahner/go-myspace/p2p"

	"github.com/bahner/go-myspace/config"

	"github.com/bahner/go-myspace/app"
)

func main() {

	ctx := context.Background()

	// Init config and common services
	config.Init(ctx)

	log := config.Log

	// Start background services
	wg := &sync.WaitGroup{}
	wg.Add(2)

	p2p.StartPubSubService(ctx, wg)
	app.StartApplication(ctx, wg)

	wg.Wait()

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
