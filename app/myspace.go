package app

// This genserver is started by the Myspace
// It is used to add children, ie. topics to the
// Myspace supervision tree

import (
	"context"
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"github.com/bahner/go-myspace/topic"
)

type Myspace struct {
	gen.Server
	ctx context.Context
}

func createMyspace(ctx context.Context) gen.ServerBehavior {

	return &Myspace{
		ctx: ctx,
	}

}

func (gr *Myspace) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	log.Infof("Initializing %s GenServer", appName)

	return nil

}

func (gr *Myspace) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Infof("Creating new topic subscription with no reply: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Myspace) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	log.Debugf("Creating new topic with a reply: %s\n", message)

	t, err := extractTopic(message)
	if err != nil {
		return nil, gen.ServerStatusIgnore
	}

	log.Debugf("Extracted topic from message: %s\n", t)

	subscribeTopic(gr.ctx, t)

	msg := etf.Tuple{
		etf.Atom("go_myspace_created_topic"),
		etf.String(t)}

	return msg, gen.ServerStatusOK
}

func (gr *Myspace) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func subscribeTopic(ctx context.Context, topicID string) {

	log.Debugf("Subscribing to topic: %s", topicID)

	process, err := n.Spawn(topicID, gen.ProcessOptions{}, topic.CreateTopicSubscription(ctx, topicID), topicID)
	if err != nil {
		switch err.Error() {
		case "resource is taken":
			log.Infof("Already subscribed to topic %s.", topicID)
		default:
			log.Errorf("Error subscribing to topic: %s", err)
		}
		return
	}
	n.RegisterName(topicID, process.Self())
}

func extractTopic(term etf.Term) (string, error) {
	log.Debugf("Extracting topic from term: %s", term)
	switch v := term.(type) {
	case []uint8:
		return string(v), nil
	case string:
		return v, nil
	case etf.Atom:
		return string(v), nil
	default:
		msg := fmt.Errorf("unexpected message type: %T", v)
		log.Error(msg)
		return "", error(msg)
	}
}
