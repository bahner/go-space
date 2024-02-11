package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
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
func getOrCreateEntity(id string) (*entity.Entity, error) {

	_e, ok := entities.Get(id)
	if ok {
		if _e, ok := _e.(*entity.Entity); ok {
			log.Debugf("found topic: %s in entities cache.", id)
			return _e, nil
		}
	}

	k, err := set.GetOrCreate(id)
	if err != nil {
		return nil, err
	}

	// // We need to publish the identity to the network, before we can create the entity
	// err = publishIdentityFromKeyset(k)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to publish identity: %w", err)
	// }

	log.Debugf("Publishing identity for entity in the background: %s", id)
	go publishIdentityFromKeyset(k)

	// I know this smells, but we just need to publish  for us to be happy.
	// We needn't wait until the rest of the world knows, to be able to send messages to the entity
	time.Sleep(3 * time.Second)

	e, err := entity.NewFromKeyset(k, id)
	if err != nil {
		return nil, err
	}

	entities.Set(id, e)

	return e, nil
}

func (e *entityCache) Set(key string, value interface{}) {
	e.store.Store(key, value)
}

func (e *entityCache) Get(key string) (interface{}, bool) {
	return e.store.Load(key)
}

func publishIdentityFromKeyset(k *set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to create DOC: %v", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to get verification method: %v", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	o := doc.DefaultPublishOptions()
	o.Force = true
	_, err = d.Publish(o)
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to publish DOC: %v", err)

	}
	log.Debugf("config.publishIdentity: published DOC: %s", d.ID)

	return nil
}
