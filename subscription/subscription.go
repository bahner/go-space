package subscription

import (
	"context"
	"errors"
	"fmt"

	"github.com/bahner/go-myspace/topic"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

type Subscription struct {
	gen.Server
	topic *topic.Topic
	ctx   context.Context
	owner gen.ProcessID
}

func New(ctx context.Context, id string) gen.ServerBehavior {

	log.Debugf("Creating new topic subscription: %s", id)

	topic, err := topic.New(id)
	if err != nil {
		fmt.Printf("Error creating topic: %s\n", err)
		return nil
	}

	log.Debugf("Created topic: %s", topic.TopicID)

	owner := createOwnerProcessId(id)
	log.Debugf("Created owner process id: %s", owner)

	return &Subscription{
		topic: topic,
		ctx:   ctx,
		owner: owner,
	}
}

func (gr *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	topic_id := gr.topic.TopicID

	log.Infof("Topic server subscribing to topic: %s\n", topic_id)
	go subscribeTopic(sp, gr) // <-- Error is here. Subscription is not working.

	return nil
}

func (gr *Subscription) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Subscription) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {

	log.Debugf("Received message: %s from: %v\n", message, from)

	action, data, err := extractActionData(message)
	if err != nil {
		log.Errorf("Error extracting action data: %s", err)
		return etf.Atom("error"), gen.ServerStatusOK
	}

	switch action {
	case "publish":
		log.Debugf("Received publish message: %s\n", data)
		gr.topic.PubSubTopic.Publish(context.Background(), data[0].([]byte))
		return etf.Atom("ok"), gen.ServerStatusOK
	case "list_peers":
		log.Debug("Received list_peers message.")
		result := gr.topic.PubSubTopic.ListPeers()
		return result, gen.ServerStatusOK
	case "get_topics":
		log.Debug("Received get_topics message.")
		result := ps.Sub.GetTopics()
		return result, gen.ServerStatusOK
	default:
		log.Debugf("Received unknown message: %s\n", data)
		return "error", gen.ServerStatusOK
	}

}

func (gr *Subscription) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func subscribeTopic(to *gen.ServerProcess, s *Subscription) {

	var sid = s.topic.TopicID
	var ctx = s.ctx

	log.Infof("Starting to listen for messages on topic: %s", sid)

	sub, err := s.topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Errorf("Error subscribing to topic: %s", sid)
	}

	log.Infof("Subscribed to topic: %s\n", sid)

	for {
		log.Debugf("Waiting for next message in topic: %s\n", sid)
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Errorf("Error getting next message: %v", err)
			continue
		}
		log.Debugf("Received message: %s\n", msg.GetData())
		sendMessage(to, s.owner, msg.GetData())
	}
}

func sendMessage(process *gen.ServerProcess, dst gen.ProcessID, data []byte) error {

	log.Debugf("Sending message to: %s", dst)

	err := process.Process.Send(dst, etf.Term(data))
	if err != nil {
		log.Error(err)
	}

	return nil
}

func createOwnerProcessId(id string) gen.ProcessID {
	return gen.ProcessID{
		Name: id,
		Node: myspaceNodeName,
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
