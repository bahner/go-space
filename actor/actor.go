package actor

import (
	"fmt"

	"github.com/bahner/go-ma/entity"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
	"github.com/bahner/go-ma/msg/envelope"

	log "github.com/sirupsen/logrus"
)

type Actor struct {
	Entity *entity.Entity

	// Incoming messages from the actor to AssertionMethod topic. It's bascially a broadcast channel.
	// But you could use it to send messages to a specific actor or to all actors in a group.
	// This is a public channel. There will need to be some generic To (recipients) in the mesage
	// for example "broadcast", so that one actor can send a message to everybody in the room.
	// That is a TODO.
	Messages chan *msg.Message

	// Incoming envelopes from the subscription to KeyAgreement topic
	// Since they are encrypted this is a private channel.
	// The envelopes are decrypted and then messages are sent to the Messages channel.
	Envelopes chan *envelope.Envelope

	// REPL channel for sending commands to the actor.
	// Messages are probably sent here from the Messages channel, after verification
	REPL chan string
}

func NewFromKeyset(k *set.Keyset, forcePublish bool) (*Actor, error) {

	a := &Actor{}
	var err error

	log.Debugf("Setting Actor Entity: %v", k)
	a.Entity, err = entity.NewFromKeyset(*k)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create Entity: %v", err)
	}

	err = a.Entity.Publish(forcePublish)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to publish Entity: %v", err)
	}
	// // We can now
	// log.Debugf("new_actor: Joining to topic: %s", room)
	// recvTopic, err := ps.Sub.Join(room)
	// if err != nil {
	// 	return nil, fmt.Errorf("new_actor: Failed to join topic: %v", err)
	// }

	// log.Debugf("new_actor: Subscribing to topic: %s", room)
	// a.From, err = recvTopic.Subscribe()
	// if err != nil {
	// 	return nil, fmt.Errorf("new_actor: Failed to subscribe to topic: %v", err)
	// }

	log.Debugf("new_actor: Actor initialized: %s", a.DID.Fragment)
	return a, nil
}
