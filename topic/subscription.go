package topic

import (
	"context"
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

type Subscription struct {
	gen.Server
	topic Topic
	ctx   context.Context
	owner gen.ProcessID
}

func CreateTopicSubscription(ctx context.Context, id string) gen.ServerBehavior {

	log.Debugf("Creating new topic subscription: %s", id)

	topic, err := getOrCreateTopic(id)
	if err != nil {
		fmt.Printf("Error creating topic: %s\n", err)
		return nil
	}

	log.Debugf("Created topic: %s", topic.TopicID)

	owner := createOwnerProcessId(id)
	log.Debugf("Created owner process id: %s", owner)

	return &Subscription{
		topic: *topic,
		ctx:   ctx,
		owner: owner,
	}
}

func (gr *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	topic_id := gr.topic.TopicID

	log.Infof("Topic server subscribing to topic: %s\n", topic_id)
	go subscribeTopic(sp, gr)

	return nil
}

func (gr *Subscription) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Subscription) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	debugstring := fmt.Sprintf("Received message: %s from: %v\n", message, from)
	fmt.Print(debugstring)
	return "ok, got it!", gen.ServerStatusOK
}

func (gr *Subscription) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func subscribeTopic(to *gen.ServerProcess, s *Subscription) {

	var sid = s.topic.TopicID

	log.Infof("Starting to listen for messages on topic: %s\n", sid)

	sub, err := s.topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Errorf("Error subscribing to topic: %s\n", sid)
	}
	defer sub.Cancel()

	log.Infof("Subscribed to topic: %s\n", sid)

	for {
		msg, err := sub.Next(s.ctx)
		if err != nil {
			log.Printf("Error getting next message: %v", err)
			continue
		}

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
