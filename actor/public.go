package actor

import (
	"fmt"

	"github.com/bahner/go-ma/msg"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (a *Actor) receivePublicMessages(sub *pubsub.Subscription) (*msg.Message, error) {

	msgData, err := sub.Next(a.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from inbox: %v", err)
	}

	m, err := msg.Unpack(string(msgData.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to unpack message: %v", err)
	}

	// Quickly check if message is signed before we try to verify it.
	// This is where DOS'ing might happen, so...
	if m.Signature == "" {
		return nil, fmt.Errorf("message has no signature")
	}

	_, err = m.Verify()
	if err != nil {
		return nil, fmt.Errorf("failed to verify message: %v", err)
	}

	return m, nil
}

func (a *Actor) handlePublicMessages(sub *pubsub.Subscription) {
	for {
		select {
		case <-a.Ctx.Done():
			// Exit goroutine when context is cancelled
			return
		default:
			// Read message from Inbox subscription
			if msg, err := a.receivePublicMessages(sub); err == nil {
				a.Messages <- msg
			}
		}
	}
}
