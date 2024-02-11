package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func (s *Subscription) handleEnvelopesLoop(ctx context.Context) {

	log.Debugf("Starting subscription envelope handling loop for topic: %s", s.entity.Topic.String())
	log.Debugf("Reading envelopes from: %v", s.entity.Envelopes)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Context for %s closed.", s.entity.Topic.String())
			return
		default:
			envelope, ok := <-s.entity.Envelopes
			if !ok {
				log.Infof("subscriptionHandleEnvelopesLoop: Envelopes channel for %s closed.", s.entity.Topic.String())
				return
			}
			log.Debugf("subscriptionHandleEnvelopesLoop: Received envelope: %s", envelope)
			msg, err := envelope.Open(s.entity.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("subscriptionHandleEnvelopesLoop: Error opening envelope: %s", err)
				continue
			}

			if msg.Verify() != nil {
				log.Errorf("subscriptionHandleEnvelopesLoop: Message not verified: %v", msg)
				continue
			}

			log.Debugf("subscriptionHandleEnvelopesLoop: open envelope and found: %v", msg)
			s.entity.Messages <- msg
		}
	}
}
