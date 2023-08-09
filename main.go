package main

import (
	"context"

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

	// Start p2p node and services
	go p2p.StartPubSubService(ctx)

	// Start Erlang node and application
	n := app.StartApplication(ctx)

	status := n.IsAlive()
	log.Printf("Node is alive: %v\n", status)

	stats := n.Stats()
	log.Printf("Node stats: %v\n", stats)

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
