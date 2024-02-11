package main

// This genserver is started by the SPACE
// It is used to add children, ie. topics to the
// SPACE supervision tree

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	log "github.com/sirupsen/logrus"
)

type SPACE struct {
	gen.Server
}

func (gr *SPACE) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	log.Infof("Initializing %s GenServer", NAME)

	return nil

}

func (gr *SPACE) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {

	log.Infof("Creating new topic subscription with no reply: %s", message)
	return gen.ServerStatusOK
}

func (gr *SPACE) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {

	log.Debugf("Creating new topic with a reply: %s", message)

	t, err := extractTopic(message)
	if err != nil {
		return nil, gen.ServerStatusIgnore
	}

	log.Debugf("Extracted topic from message: %s", t)

	go subscribeTopic(t)

	msg := etf.Tuple{
		etf.Atom("go_space_creating_topic"),
		etf.String(t)}

	return msg, gen.ServerStatusOK
}

func (gr *SPACE) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s", message)
	return gen.ServerStatusOK
}

func subscribeTopic(topicID string) {

	log.Debugf("Spawning subscription process: %s", topicID)

	n = getNode()
	log.Debugf("Node name: %s is alive. %t", n.Name(), n.IsAlive())

	log.Debugf("Subscribing to topic: %s", topicID)

	sub := New(topicID)

	log.Debugf("Subscription: %s", sub)

	process, err := n.Spawn(topicID, gen.ProcessOptions{}, sub, topicID)
	if err != nil {
		switch err.Error() {
		case "resource is taken":
			log.Infof("Already subscribed to topic %s.", topicID)
		default:
			log.Errorf("Error subscribing to topic: %s", err)
		}
		return
	}
	log.Debugf("Registering process: %s with name: %s", process.Self(), topicID)
	n.RegisterName(topicID, process.Self())
}

func extractTopic(term etf.Term) (string, error) {

	log.Debugf("Extracting topic from term: %s", term)
	switch v := term.(type) {
	case []uint8:
		log.Debugf("Converting []uint8 to string: %s", string(v))
		return string(v), nil
	case string:
		log.Debugf("Converting string to string: %s", v)
		return v, nil
	case etf.Atom:
		log.Debugf("Converting atom to string: %s", string(v))
		return string(v), nil
	default:
		msg := fmt.Errorf("unexpected message type: %T", v)
		log.Error(msg)
		return "", error(msg)
	}
}
