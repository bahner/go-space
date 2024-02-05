package topic

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var (
	topics sync.Map
)

type Topic struct {
	PubSubTopic *p2ppubsub.Topic
	TopicID     string
}

func New(topicID string) (*Topic, error) {

	service := pubsub.Get()
	if service == nil {
		return nil, fmt.Errorf("topic: error getting pubsub service")
	}

	log.Debugf("Looking for topic: %s in topics map", topicID)
	topic, ok := topics.Load(topicID)
	if ok {
		if t, ok := topic.(*Topic); ok {
			log.Debugf("Found topic: %s in topics map", topicID)
			return t, nil
		}
	}

	log.Debugf("Topic: %s not found in topics map, creating new topic", topicID)
	pubSubTopic, err := service.Join(topicID)
	if err != nil {
		log.Errorf("Error joining topic: %s", err)
		return nil, err
	}

	topic = &Topic{
		PubSubTopic: pubSubTopic,
		TopicID:     topicID,
	}

	topics.Store(topicID, topic)

	return topic.(*Topic), nil
}
