package app

// This genserver is started by the Space
// It is used to add children, ie. topics to the
// Space supervision tree

import (
	"context"
	"fmt"

	"github.com/bahner/go-space/config"
	"github.com/bahner/go-space/subscription"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	log "github.com/sirupsen/logrus"
)

type Space struct {
	gen.Server
	ctx context.Context
}

func createSpace(ctx context.Context) gen.ServerBehavior {

	return &Space{
		ctx: ctx,
	}

}

func (gr *Space) Init(sp *gen.ServerProcess, args ...etf.Term) error {
	appName := config.AppName

	log.Infof("Initializing %s GenServer", appName)

	return nil

}

func (gr *Space) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {

	log.Infof("Creating new topic subscription with no reply: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Space) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {

	log.Debugf("Creating new topic with a reply: %s\n", message)

	t, err := extractTopic(message)
	if err != nil {
		return nil, gen.ServerStatusIgnore
	}

	log.Debugf("Extracted topic from message: %s\n", t)

	subscribeTopic(gr.ctx, t)

	msg := etf.Tuple{
		etf.Atom("go_space_created_topic"),
		etf.String(t)}

	return msg, gen.ServerStatusOK
}

func (gr *Space) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Debugf("Received message: %s\n", message)
	return gen.ServerStatusOK
}

func subscribeTopic(ctx context.Context, topicID string) {

	n = getNode()
	log.Debugf("Node name: %s is alive. %t", n.Name(), n.IsAlive())

	log.Debugf("Subscribing to topic: %s", topicID)

	sub := subscription.New(ctx, topicID)
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
