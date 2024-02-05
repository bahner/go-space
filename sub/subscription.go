package sub

import (
	"context"
	"errors"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Subscription struct {
	gen.Server
	Cancel context.CancelFunc
	Ctx    context.Context
	topic  *p2ppubsub.Topic
	owner  gen.ProcessID
	sp     *gen.ServerProcess
	sub    *p2ppubsub.Subscription
}

func New(id string) gen.ServerBehavior {

	log.Debugf("Creating new topic subscription: %s", id)

	topic, err := getOrCreateTopic(id)
	if err != nil {
		log.Errorf("Error creating topic: %v", err)
		return nil
	}

	log.Debugf("Created topic: %s", topic.String())

	owner := createOwnerProcessId(id)
	log.Debugf("Created owner process id: %s", owner)

	return &Subscription{
		topic: topic,
		owner: owner,
	}
}

func (s *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	s.sp = sp

	log.Infof("Subscription init subscribing to topic: %s", s.topic.String())
	go s.subscriptionLoop() // <-- Error is here. Subscription is not working.

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
		s.topic.Publish(context.Background(), data[0].([]byte))
		return etf.Atom("ok"), gen.ServerStatusOK
	case "list_peers":
		log.Debug("Received list_peers message.")
		result := s.topic.ListPeers()
		return result, gen.ServerStatusOK
	case "get_topics":
		log.Debug("Received get_topics message.")
		ps := pubsub.Get()
		if ps == nil {
			log.Errorf("Error getting pubsub service")
			return "error", gen.ServerStatusOK
		}
		result := ps.GetTopics()
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

	var sid = s.topic.String()

	log.Infof("Starting to listen for messages on topic: %s", sid)

	sub, err := s.topic.Subscribe()
	if err != nil {
		log.Errorf("Error subscribing to topic: %s", sid)
	}

	log.Infof("Subscribed to topic: %s", sid)

	for {
		log.Debugf("Waiting for next message in topic: %s", sid)
		msg, err := sub.Next(context.Background())
		if err != nil {
			log.Errorf("Error getting next message: %v", err)
			continue
		}
		log.Debugf("Received message: %s", msg.GetData())
		s.sendMessageToOwner(msg.GetData())
	}
}

func (s *Subscription) sendMessageToOwner(data []byte) error {
	log.Debugf("Sending message to: %s", s.owner)

	err := s.sp.Process.Send(s.owner, etf.Term(data))
	if err != nil {
		log.Error(err)
	}

	return nil
}

func createOwnerProcessId(id string) gen.ProcessID {
	return gen.ProcessID{
		Name: id,
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
