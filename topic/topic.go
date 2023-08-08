package topic

import (
	"myspace-pubsub/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Topic struct {
	PubSubTopic *pubsub.Topic
	TopicID     string
}

func getOrCreateTopic(topicID string) (*Topic, error) {

	ps := p2p.PubSubService

	log.Debugf("Looking for topic: %s in topics map", topicID)
	topic, ok := topics.Load(topicID)
	if ok {
		if t, ok := topic.(*Topic); ok {
			log.Debugf("Found topic: %s in topics map", topicID)
			return t, nil
		}
	}

	log.Debugf("Topic: %s not found in topics map, creating new topic", topicID)
	log.Debugf("Joining topic: %s", topicID)
	pubSubTopic, err := ps.Join(topicID)
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
