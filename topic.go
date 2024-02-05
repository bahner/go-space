package main

import (
	"fmt"
	"sync"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var topics *cache

type cache struct {
	store sync.Map
}

func init() {
	topics = new(cache)
}

func getOrCreateTopic(topicID string) (*p2ppubsub.Topic, error) {

	topic, ok := topics.Get(topicID)
	if ok {
		if t, ok := topic.(*p2ppubsub.Topic); ok {
			log.Debugf("Found topic: %s in topics cache.", topicID)
			return t, nil
		}
	}

	if p.PubSub == nil {
		return nil, fmt.Errorf("pubsub service not available")
	}

	log.Debugf("Topic: %s not found in topics map, creating new topic", topicID)
	t, err := p.PubSub.Join(topicID)
	if err != nil {
		log.Errorf("Error joining topic: %s", err)
		return nil, err
	}

	topics.Set(topicID, t)

	return t, nil
}

func (t *cache) Set(key string, value interface{}) {
	t.store.Store(key, value)
}

func (t *cache) Get(key string) (interface{}, bool) {
	return t.store.Load(key)
}
