package main

import (
	"context"
	"encoding/json"

	"github.com/ergo-services/ergo/etf"
	log "github.com/sirupsen/logrus"
)

func (s *Subscription) handleMessagesLoop(ctx context.Context) {

	t := s.actor.Entity.Topic.String()
	messages := s.actor.Entity.Messages

	log.Debugf("Starting subscription message handling loop for topic: %s", t)
	log.Debugf("Reading messages from: %v", messages)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Context for %s closed.", t)
			return
		default:
			message, ok := <-messages
			if !ok {
				log.Infof("Messages channel for %s closed.", t)
				return
			}

			log.Debugf("subscriptionHandleMessagesLoop: received message: %v", message)

			// Marshal the message and send it to the owner
			msgJson, err := json.Marshal(message)
			if err != nil {
				log.Errorf("Error marshaling message: %s", err)
				continue
			}

			// Send message as JSON to owner
			s.deliverMessage(msgJson)
		}
	}
}

func (s *Subscription) deliverMessage(data []byte) error {

	log.Debugf("Delivering message: %s to owner: %s", data, s.owner)
	err := s.sp.Process.Send(s.owner, etf.Term(data))
	if err != nil {
		log.Error(err)
	}

	return nil
}
