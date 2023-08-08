package main

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	pubsub "myspace/pubsub"
)

var ps = pubsub.PubSubService

type Room struct {
	gen.Server
	topic Topic
}

func createRoom(id string) gen.ServerBehavior {

	topic, err := getOrCreateTopic(id)
	if err != nil {
		fmt.Printf("Error creating topic: %s\n", err)
		return nil
	}

	return &Room{
		topic: *topic,
	}
}

func (gr *Room) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	topic_id := gr.topic.TopicID

	fmt.Printf("Initializing Go Room Process with room ID: %s\n", topic_id)

	_, err := ps.Join(topic_id)
	if err != nil {
		fmt.Printf("Error joining topic: %s\n", err)
		return nil
	}
	fmt.Printf("Joined topic: %s\n", topic_id)

	fmt.Printf("Starting to listen for messages on topic: %s\n", topic_id)
	// go subscribeTopic(p.Self(), gr.topic)
	go subscribeTopic(sp, gr.topic)

	return nil
}

func (gr *Room) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Room) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	debugstring := fmt.Sprintf("Received message: %s from: %v\n", message, from)
	fmt.Print(debugstring)
	return "ok, got it!", gen.ServerStatusOK
}

func (gr *Room) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("Received message: %s\n", message)
	return gen.ServerStatusOK
}
