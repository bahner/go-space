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
}

func CreateTopicSubscription(ctx context.Context, id string) gen.ServerBehavior {

	topic, err := getOrCreateTopic(id)
	if err != nil {
		fmt.Printf("Error creating topic: %s\n", err)
		return nil
	}

	return &Subscription{
		topic: *topic,
		ctx:   ctx,
	}
}

func (gr *Subscription) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	topic_id := gr.topic.TopicID

	fmt.Printf("Initializing Go Subscription Process with Subscription ID: %s\n", topic_id)

	_, err := ps.Join(topic_id)
	if err != nil {
		fmt.Printf("Error joining topic: %s\n", err)
		return nil
	}
	fmt.Printf("Joined topic: %s\n", topic_id)

	fmt.Printf("Starting to listen for messages on topic: %s\n", topic_id)
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

	log.Infof("Starting to listen for messages on topic: %s\n", s.topic.TopicID)

	sub, err := s.topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()

	for {
		msg, err := sub.Next(s.ctx)
		if err != nil {
			log.Printf("Error getting next message: %v", err)
			continue
		}

		// Send the received message back to the GenServer
		to.Send(to.Self(), etf.Term(string(msg.GetData())))
	}
}
