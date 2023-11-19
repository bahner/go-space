package actor

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/entity"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	log "github.com/sirupsen/logrus"
)

const MESSAGES_BUFFERSIZE = 100

type Actor struct {

	// This context is used to cancel the Listen() function.
	Ctx context.Context

	// All actors must be entities.
	// Ideally they should be the same, but then ma becomes a bit too opinionated.
	Entity *entity.Entity

	// Private is the topic where we receive envelopes from other actors.
	// It's basically a private channel with the DIDDocument keyAgreement as topic.
	Private *pubsub.Topic

	// We basically receive signed messages from the room we're in here.
	// It's basically a public channel with the assertionMethod from the DIDDocument of
	// the room we're in as topic.
	// Others can subscribe to this topic and send us messages, as long as they are signed.
	Public *pubsub.Topic

	// Incoming messages from the actor to AssertionMethod topic. It's bascially a broadcast channel.
	// But you could use it to send messages to a specific actor or to all actors in a group.
	// This is a public channel. There will need to be some generic To (recipients) in the mesage
	// for example "broadcast", so that one actor can send a message to everybody in the room.
	// That is a TODO.
	// We receive the message contents here after verification or decryption.
	Messages chan *msg.Message
}

// Creates a new actor from an entity.
// Takes a pubsub.PubSub service, an entity and a forcePublish flag.
// The forcePublish is to override existing keys in IPFS.
func New(ctx context.Context, ps *pubsub.PubSub, e *entity.Entity, forcePublish bool) (*Actor, error) {

	var err error
	a := &Actor{}

	// Firstly create assign entity to actor
	a.Entity = e

	// Create topic for incoming envelopes
	a.Private, err = ps.Join(a.Entity.Doc.KeyAgreement)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to join topic: %v", err)
	}

	// Create subscription to topic for incoming messages
	a.Public, err = ps.Join(a.Entity.Doc.AssertionMethod)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to join topic: %v", err)
	}

	// Set the messages channel
	a.Messages = make(chan *msg.Message, MESSAGES_BUFFERSIZE)

	// Publish the entity
	err = a.Entity.Publish(forcePublish)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to publish Entity: %v", err)
	}

	log.Debugf("new_actor: Actor initialized: %s", a.Entity.DID.Fragment)
	return a, nil

}

// Creates a new actor from a keyset.
// Takes a pubsub.PubSub service, a keyset and a forcePublish flag.
func NewFromKeyset(ctx context.Context, ps *pubsub.PubSub, k *set.Keyset, forcePublish bool) (*Actor, error) {

	log.Debugf("Setting Actor Entity: %v", k)
	e, err := entity.NewFromKeyset(*k)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create Entity: %v", err)
	}

	return New(ctx, ps, e, forcePublish)
}
