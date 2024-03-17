package main

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
)

var entities *entityCache

type entityCache struct {
	store sync.Map
}

func init() {
	entities = new(entityCache)
}

// GetOrCreateEntity returns an entity from the cache or creates a new one
// The id is just the uniqgue name of the calling entity, not the full DID
func getOrCreateEntity(id string) (*actor.Actor, error) {

	// Attempt to retrieve the entity from cache.
	// This is runtime, so entities will be generated at least once.
	if cachedEntity, ok := entities.Get(id); ok {
		if entity, ok := cachedEntity.(*actor.Actor); ok {
			log.Debugf("found topic: %s in entities cache.", id)
			return entity, nil // Successfully type-asserted and returned
		}
	}

	// Entity not found in cache, proceed to creation
	log.Debugf("getOrCreateEntity: GetOrCreateKeyset from vault: %s", id)
	k, err := getOrCreateKeysetFromVault(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create keyset: %w", err)
	}

	// Assuming entity.NewFromKeyset returns *actor.Actor
	e, err := actor.NewFromKeyset(k)
	if err != nil {
		return nil, fmt.Errorf("failed to create entity: %w", err)
	}

	e.Entity.Doc, err = e.CreateDocument(e.Entity.DID.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to create DID Document: %w", err)
	}

	// Force publication of document.
	o := doc.DefaultPublishOptions()
	o.Force = true
	e.Entity.Doc.Publish(o)

	// Cache the newly created entity for future retrievals
	entities.Set(id, e)

	return e, nil
}

func (e *entityCache) Set(key string, value interface{}) {
	e.store.Store(key, value)
}

func (e *entityCache) Get(key string) (interface{}, bool) {
	return e.store.Load(key)
}
