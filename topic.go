package main

import (
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var (
	pub    *pubsub.PubSub
	topics sync.Map
)

type Topic struct {
	PubSubTopic *pubsub.Topic
	TopicID     string
}

func getOrCreateTopic(topicID string) (*Topic, error) {
	topic, ok := topics.Load(topicID)
	if ok {
		if t, ok := topic.(*Topic); ok {
			return t, nil
		}
	}

	pubSubTopic, err := pub.Join(topicID)
	if err != nil {
		return nil, err
	}

	topic = &Topic{
		PubSubTopic: pubSubTopic,
		TopicID:     topicID,
	}

	topics.Store(topicID, topic)

	return topic.(*Topic), nil
}
