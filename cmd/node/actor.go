package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	log "github.com/sirupsen/logrus"
)

func initActorOrPanic() *actor.Actor {
	// The actor is needed for initialisation of the WebHandler.
	fmt.Println("Creating actor from keyset...")
	a, err := actor.NewFromKeyset(config.ActorKeyset())
	if err != nil {
		log.Debugf("error creating actor: %s", err)
	}

	id := a.Entity.DID.Id

	fmt.Println("Creating and setting DID Document for actor...")
	err = a.CreateAndSetDocument(id)
	if err != nil {
		panic(fmt.Sprintf("error creating document: %s", err))
	}

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a == nil || a.Verify() != nil {
		panic(fmt.Sprintf("%s is not a valid actor: %v", id, err))
	}

	_, err = entity.GetOrCreateFromDID(a.Entity.DID, false)
	if err != nil {
		panic(fmt.Sprintf("error getting or creating entity: %s", err))
	}

	return a
}
