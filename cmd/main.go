package main

import (
	"context"

	"github.com/bahner/go-space/app"
	"github.com/bahner/go-space/global"
	"github.com/bahner/go-space/keeper"

	"github.com/bahner/go-space/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	// Init config and logger
	config.Init(ctx)

	// Start background services
	global.InitAndStartServices(ctx)

	// Start application
	app.StartApplication(ctx)

	_secret := []byte("secret")

	k := keeper.New("space")
	defer k.Close()

	safe_secret, err := keeper.Encrypt(k, _secret)
	if err != nil {
		log.Fatal(err)
	}
	log.Warnf("safe_secret: %v", safe_secret)

	select {}
}
