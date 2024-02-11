package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	MESSAGE_CHANNEL_SIZE  = 100
	ENVELOPE_CHANNEL_SIZE = 100
)

var (
	subscriptionsCache = sync.Map{}
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

	if s, ok := subscriptionsCache.Load(id); ok {
		log.Debugf("Found existing genServer: %s", id)
		return s.(gen.ServerBehavior)
	}

	log.Debugf("Creating new genServer: %s", id)

	entity, err := getOrCreateEntity(id)
	if err != nil {
		log.Errorf("Error getting or creating entity: %s", err)
		return nil
	}

	log.Debugf("Created topic: %s", entity.Topic.String())

	// The owner is identified by the fragment of the DID
	// It's the local name ad ID of the owner of the entity
	owner := createOwnerProcessId(did.GetFragment(id))
	log.Debugf("Created owner process id: %s", owner)

	s := &Subscription{
		owner:     owner,
		entity:    entity,
		messages:  make(chan *msg.Message, MESSAGE_CHANNEL_SIZE),
		envelopes: make(chan *msg.Envelope, ENVELOPE_CHANNEL_SIZE),
	}

	subscriptionsCache.Store(id, s)

	return s
}

func (s *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	s.sp = sp // Save the server process, so we can send messages from it

	ctx := context.Background()

	log.Infof("Subscription init subscribing to topic: %s", s.entity.DID.String())

	log.Debugf("Subscription entity: %v", s.entity)
	go s.entity.Subscribe(ctx, s.entity)
	go s.subscribe()

	sp.Process.Send(s.owner, etf.Tuple{
		etf.Atom(":go_space_topic_subscription_created"),
		etf.String(s.entity.Topic.String()),
	})

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
		s.entity.Topic.Publish(s.entity.Ctx, data[0].([]byte))
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

func (s *Subscription) Terminate(sp *gen.ServerProcess, reason string) {

	// Close the topic.
	s.entity.Cancel()

	sp.Kill()

	log.Debugf("Terminating subscription: %s", reason)
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

func (s *Subscription) subscribe() {
	log.Debugf("Starting subscription loop: %s", s.entity.Topic.String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.handleMessagesLoop(ctx)
	go s.handleEnvelopesLoop(ctx)

	<-ctx.Done()
	log.Infof("Context for %s closed.", s.entity.Topic.String())
}
