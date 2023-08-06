package main

import (
	"fmt"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
)

func (ga *Avatar) Init(p *ergo.Process, args ...interface{}) (state interface{}) {
	fmt.Printf("Initializing Go Avatar Process with avatar ID: %s\n", args[0])

	// Create a new GossipSub topic for this avatar
	topic, _ := ps.Join(args[0].(string))

	ga.topic = topic

	return AvatarState{avatarID: args[0].(string)}
}

func (ga *Avatar) HandleCast(message etf.Term, state interface{}) (string, interface{}) {
	fmt.Printf("Received message: %s\n", message)
	return "noreply", state
}

func (ga *Avatar) HandleCall(from ergo.ProcessFrom, message etf.Term, state interface{}) (string, etf.Term, interface{}) {
	fmt.Printf("Received message: %s from: %v\n", message, from)
	return "reply", etf.Term("ok"), state
}
