package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	MESSAGE_CHANNEL_SIZE  = 100
	ENVELOPE_CHANNEL_SIZE = 100
)

type Subscription struct {
	gen.Server
	sp     *gen.ServerProcess
	owner  gen.ProcessID
	entity *entity.Entity

	messages  chan *msg.Message
	envelopes chan *msg.Envelope
}

func (s *Subscription) Verify() error {
	if s.owner.Node == "" {
		return fmt.Errorf("owner node is empty")
	}
	if s.owner.Name == "" {
		return fmt.Errorf("owner name is empty")
	}
	if s.entity == nil {
		return fmt.Errorf("entity is nil")
	}
	if s.messages == nil {
		return fmt.Errorf("messages channel is nil")
	}
	if s.envelopes == nil {
		return fmt.Errorf("envelopes channel is nil")
	}
	return nil
}

func (s *Subscription) IsValid() bool {
	return s.Verify() == nil
}

func New(id string) gen.ServerBehavior {

	log.Debugf("Creating new genServer: %s", id)

	entity, err := getOrCreateEntity(id)
	if err != nil {
		log.Errorf("Error creating entity: %v", err)
		return nil
	}

	log.Debugf("Created topic: %s", entity.Topic.String())

	// The owner is identified by the fragment of the DID
	// It's the local name ad ID of the owner of the entity
	owner := createOwnerProcessId(did.GetFragment(id))
	log.Debugf("Created owner process id: %s", owner)

	return &Subscription{
		owner:     owner,
		entity:    entity,
		messages:  make(chan *msg.Message, MESSAGE_CHANNEL_SIZE),
		envelopes: make(chan *msg.Envelope, ENVELOPE_CHANNEL_SIZE),
	}
}

func (s *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	s.sp = sp // Save the server process, so we can send messages from it

	log.Infof("Subscription init subscribing to topic: %s", s.entity.DID.String())
	go s.subscriptionLoop() // <-- Error is here. Subscription is not working.

	sp.Process.Send(s.owner, etf.Atom(":go_space_topic_created"))

	return nil
}

func (s *Subscription) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s", message)
	return gen.ServerStatusOK
}

func (s *Subscription) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {

	log.Debugf("Received message: %s from: %v", message, from)

	action, data, err := extractActionData(message)
	if err != nil {
		log.Errorf("Error extracting action data: %s", err)
		return etf.Atom("error"), gen.ServerStatusOK
	}

	switch action {

	case "publish":
		log.Debugf("Received publish message: %s", data)
		s.entity.Topic.Publish(context.Background(), data[0].([]byte))
		return etf.Atom("ok"), gen.ServerStatusOK

	case "list_peers":
		log.Debug("Received list_peers message.")
		result := s.entity.Topic.ListPeers()
		return result, gen.ServerStatusOK

	case "get_topics":
		log.Debug("Received get_topics message.")
		result := p.PubSub.GetTopics()
		return result, gen.ServerStatusOK

	default:
		log.Debugf("Received unknown message: %s", data)
		return "error", gen.ServerStatusOK
	}

}

func (s *Subscription) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s", message)
	return gen.ServerStatusOK
}

func (s *Subscription) subscriptionLoop() {

	sub, err := s.entity.Subscribe()
	if err != nil {
		log.Errorf("Error subscribing to topic: %s", err)
		return
	}
	defer sub.Cancel()

	// Start a debug loop if log level is debug
	if log.GetLevel() == log.DebugLevel {
		go s.debugLoop()
	}

	// Start the message and envelope handling loops
	// which listen on the channels for messages and envelopes
	go s.handleEnvelopesLoop()
	go s.handleMessagesLoop()

	var t = s.entity.Topic.String()

	log.Infof("Starting to listen for messages on topic: %s", t)

	for {
		log.Debugf("Waiting for next message in topic: %s", t)
		pubsubMessage, ok := <-sub.Messages
		if !ok {
			log.Infof("Subscription: %s closed.", t)
			return
		}

		if s.handleMessage(pubsubMessage.Data) == nil {
			log.Debugf("Message: %s handled.", pubsubMessage.Data)
			continue
		}

		if s.handleEnvelope(pubsubMessage.Data) == nil {
			log.Debugf("Envelope: %s handled.", pubsubMessage.Data)
			continue
		}

		log.Errorf("Error handling message: %s", pubsubMessage.Data)

	}
}

func (s *Subscription) handleMessage(data []byte) error {

	msg := new(msg.Message)
	err := cbor.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("error unmarshaling message: %v", err)
	}

	// If the message is verified, we pass it on and continue prosessing new messages
	if msg.Verify() != nil {
		return fmt.Errorf("message not verified: %v", msg)
	}

	// Message is verified, we pass it on
	s.messages <- msg

	return nil
}

func (s *Subscription) handleEnvelope(data []byte) error {

	envelope := new(msg.Envelope)
	err := cbor.Unmarshal(data, envelope)
	if err != nil {
		return fmt.Errorf("error unmarshaling envelope: %v", err)
	}

	// Sanity check the envelope
	if envelope.EncryptedContent == nil {
		return fmt.Errorf("envelope not verified: %v", envelope)
	}

	if envelope.EncryptedHeaders == nil {
		return fmt.Errorf("envelope not verified: %v", envelope)
	}

	if envelope.EphemeralKey == nil {
		return fmt.Errorf("envelope not verified: %v", envelope)
	}

	// Seems OK, we pass it on
	s.envelopes <- envelope

	return nil
}

func (s *Subscription) deliverMessage(data []byte) error {

	log.Debugf("Delivering message: %s to owner: %s", data, s.owner)
	err := s.sp.Process.Send(s.owner, etf.Term(data))
	if err != nil {
		log.Error(err)
	}

	return nil
}

func createOwnerProcessId(id string) gen.ProcessID {

	fragment := did.GetFragment(id)

	return gen.ProcessID{
		Name: fragment,
		Node: viper.GetString("node.space"),
	}
}

func extractActionData(term etf.Term) (etf.Atom, []etf.Term, error) {
	// If the term is just an Atom
	if atom, ok := term.(etf.Atom); ok {
		return atom, nil, nil
	}

	// If the term is a Tuple
	tuple, ok := term.(etf.Tuple)
	if !ok || len(tuple) == 0 {
		return "", nil, errors.New("term is not a tuple or is empty")
	}

	command, ok := tuple[0].(etf.Atom)
	if !ok {
		return "", nil, errors.New("first element is not an atom")
	}

	// Return the command and the rest of the tuple
	return command, tuple[1:], nil
}

func (s *Subscription) Terminate(sp *gen.ServerProcess, reason string) {

	// Close the topic.
	s.entity.Subscription.Cancel()
	s.entity.Topic.Close()

	sp.Kill()

	log.Debugf("Terminating subscription: %s", reason)
}

func (s *Subscription) debugLoop() {

	for {

		if s.entity.Subscription == nil {
			log.Debugf("Subscription: %s is nil.", s.entity.Topic.String())
			return
		} else {
			log.Debugf("Subscription: %s is alive with peers: %v", s.entity.Topic.String(), s.entity.Topic.ListPeers())
		}
		time.Sleep(viper.GetDuration("node.debug_interval"))
	}
}

func (s *Subscription) handleEnvelopesLoop() {

	for {
		envelope, ok := <-s.envelopes
		if !ok {
			log.Infof("Envelopes channel for %s closed.", s.entity.Topic.String())
			return
		}
		log.Debugf("Received envelope: %s", envelope)
		msg, err := envelope.Open(s.entity.Keyset.EncryptionKey.PrivKey[:])
		if err != nil {
			log.Errorf("Error opening envelope: %s", err)
			continue
		}

		if msg.Verify() != nil {
			log.Errorf("Message not verified: %v", msg)
			continue
		}

		s.messages <- msg
	}
}

func (s *Subscription) handleMessagesLoop() {

	for {
		message, ok := <-s.messages
		if !ok {
			log.Infof("Messages channel for %s closed.", s.entity.Topic.String())
			return
		}

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
