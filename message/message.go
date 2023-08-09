package message

import (
	"time"
)

// To and From are ipns addresses. So there is no need to have a
// TopicID field. The To and From fields are the TopicIDs.

// To and From are required to publish their did to the network.
// This is so that the network can verify that the message is
// from the correct publisher. This is done by verifying the
// signature of the message.

// Hence we need to lookup the did of the publisher from the
// network. We should keep a cache of the dids of the publishers.
// If message doesn't have a valid signature, then we should
// lookup the did of the publisher from the network.

type Message struct {
	From     string
	To       string
	Data     []byte
	Created  time.Time
	Received time.Time
}

func New(from string, to string, data []byte) *Message {
	return &Message{
		From:    from, // This is actually a TopicID
		To:      to,   // This is actually a TopicID
		Data:    data,
		Created: time.Now(),
	}
}

func verifyMessage(msg *Message) bool {
	return true
}
