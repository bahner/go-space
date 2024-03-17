package main

import (
	"errors"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const name = "node"

func initConfig() {

	// Always parse the flags first
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.SetProfile(name)
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node := generateActorIdentitiesOrPanic(config.Profile())
		actorConfig := configTemplate(actor, node)
		config.Generate(actorConfig)
		os.Exit(0)
	}

	// At this point an actor *must* be initialized
	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}

func generateActorIdentitiesOrPanic(name string) (string, string) {
	actor, node, err := config.GenerateActorIdentities(name)
	if err != nil {
		if errors.Is(err, doc.ErrAlreadyPublished) {
			log.Warnf("Actor document already published: %v", err)
		} else {
			log.Fatal(err)
		}
	}
	return actor, node
}
