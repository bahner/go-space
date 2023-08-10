package topic

import (
	ps "github.com/bahner/go-myspace/p2p/pubsub"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Topic struct {
	PubSubTopic *pubsub.Topic
	TopicID     string
}

func getOrCreateTopic(topicID string) (*Topic, error) {

	service := ps.Service

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
