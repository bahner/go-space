package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func subscribeTopic(to *gen.ServerProcess, topic Topic) {

	log.Infof("Starting to listen for messages on topic: %s\n", topic.TopicID)

	sub, err := topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()

	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Printf("Error getting next message: %v", err)
			continue
		}

		// Send the received message back to the GenServer
		// gen.Process.Send(to, etf.Term(string(msg.GetData())))
		to.Send(to.Self(), etf.Term(string(msg.GetData())))
	}
}

func spawnAndRegisterRoom(roomID string) {
	process, err := n.Spawn(roomID, gen.ProcessOptions{}, createRoom(roomID), roomID)
	if err != nil {
		log.Fatal(err)
	}
	n.RegisterName(roomID, process.Self())
}
