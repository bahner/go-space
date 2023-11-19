package actor

import (
	"fmt"

	"github.com/bahner/go-ma/msg"
	"github.com/bahner/go-ma/msg/envelope"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (a *Actor) receiveFromInbox(sub *pubsub.Subscription) (*msg.Message, error) {

	msgData, err := sub.Next(a.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from inbox: %v", err)
	}

	e, err := envelope.UnmarshalFromCBOR(msgData.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal envelope from CBOR: %v", err)
	}

	message, err := e.Open(a.Entity.Keyset.EncryptionKey.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to open envelope: %v", err)
	}

	return message, nil
}

func (a *Actor) handleInboxSubscription(sub *pubsub.Subscription) {
	for {
		select {
		case <-a.Ctx.Done():
			// Exit goroutine when context is cancelled
			return
		default:
			// Read message from Inbox subscription
			if msg, err := a.receiveFromInbox(sub); err == nil {
				a.Messages <- msg
			}
		}
	}
}
