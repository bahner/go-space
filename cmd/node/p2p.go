package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
)

func initP2P() (P2P *p2p.P2P, err error) {
	fmt.Println("Initialising libp2p...")

	// Everyone needs a connection manager.
	cm, err := connmgr.Init()
	if err != nil {
		panic(fmt.Errorf("pong: failed to create connection manager: %w", err))
	}
	cg := connmgr.NewConnectionGater(cm)

	d, err := DHT(cg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize dht: %v", err))
	}
	return p2p.Init(d)
}
