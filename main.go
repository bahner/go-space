package main

import (
	"context"
	"fmt"
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type goRoom struct {
	ergo.GenServer
	process *ergo.Process
	topic   *pubsub.Topic
}

type goAvatar struct {
	ergo.GenServer
	process *ergo.Process
	topic   *pubsub.Topic
}

type goRoomState struct {
	roomID string
}

type goAvatarState struct {
	avatarID string
}

var (
	h  host.Host
	ps *pubsub.PubSub
)

func init() {
	// Create a new libp2p Host that listens on a random TCP port
	h, _ = libp2p.New(context.Background())

	// Create a new GossipSub instance
	ps, _ = pubsub.NewGossipSub(context.Background(), h)
}

func (gr *goRoom) Init(p *ergo.Process, args ...interface{}) (state interface{}) {
	fmt.Printf("Initializing Go Room Process with room ID: %s\n", args[0])

	// Create a new GossipSub topic for this room
	topic, _ := ps.Join(args[0].(string))

	gr.topic = topic

	return goRoomState{roomID: args[0].(string)}
}

func (ga *goAvatar) Init(p *ergo.Process, args ...interface{}) (state interface{}) {
	fmt.Printf("Initializing Go Avatar Process with avatar ID: %s\n", args[0])

	// Create a new GossipSub topic for this avatar
	topic, _ := ps.Join(args[0].(string))

	ga.topic = topic

	return goAvatarState{avatarID: args[0].(string)}
}

func (gr *goRoom) HandleCast(message etf.Term, state interface{}) (string, interface{}) {
	fmt.Printf("Received message: %s\n", message)
	return "noreply", state
}

func (ga *goAvatar) HandleCast(message etf.Term, state interface{}) (string, interface{}) {
	fmt.Printf("Received message: %s\n", message)
	return "noreply", state
}

func (gr *goRoom) HandleCall(from ergo.ProcessFrom, message etf.Term, state interface{}) (string, etf.Term, interface{}) {
	fmt.Printf("Received message: %s from: %v\n", message, from)
	return "reply", etf.Term("ok"), state
}

func (ga *goAvatar) HandleCall(from ergo.ProcessFrom, message etf.Term, state interface{}) (string, etf.Term, interface{}) {
	fmt.Printf("Received message: %s from: %v\n", message, from)
	return "reply", etf.Term("ok"), state
}

func main() {
	node := ergo.CreateNode("go@localhost", "secret", ergo.NodeOptions{})
	supOpts := ergo.SupervisorOptions{
		Strategy: ergo.RestartAll,
		Intensity: 1,
		Period: 5,
	}

	supSpec := ergo.SupervisorSpec{
		Name: "gameSupervisor",
		Children: []ergo.SupervisorChildSpec{
			{
				Name: "goRoom",
				ChildGenServer: ergo.SupervisorChildGenServer{
					Args: []interface{}{"room1"},
					Func: func() ergo.GenServer { return &goRoom{} },
				},
			},
			{
				Name: "goAvatar",
				ChildGenServer: ergo.SupervisorChildGenServer{
					Args: []interface{}{"avatar1"},
					Func: func() ergo.GenServer { return &goAvatar{} },
				},
			},
		},
	}

	_, _ = node.Supervisor(supOpts, supSpec)
	select {}
}
