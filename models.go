package main

import (
	"sync"

	"github.com/ergo-services/ergo"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Topic represents a chatroom topic
type Topic struct {
	Mutex       sync.Mutex
	PubSubTopic *pubsub.Topic
	TopicID     string
}

type Room struct {
	ergo.GenServer
	process *ergo.Process
	topic   *pubsub.Topic
}

type Avatar struct {
	ergo.GenServer
	process *ergo.Process
	topic   *pubsub.Topic
}

type RoomState struct {
	roomID string
}

type AvatarState struct {
	avatarID string
}
