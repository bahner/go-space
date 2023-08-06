package main

import (
	"fmt"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func createRoomSupervisor() gen.SupervisorBehavior {
	return &RoomSupervisor{}
}

type RoomSupervisor struct {
	gen.Supervisor
}

func (sup *MySup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: "mysup",
		Children: []gen.SupervisorChildSpec{
			{
				Name:  "myactor",
				Child: createMyActor(),
			},
		},
		Strategy: gen.SupervisorStrategy{
			Type:      gen.SupervisorStrategyOneForOne,
			Intensity: 2, // How big bursts of restarts you want to tolerate.
			Period:    5, // In seconds.
			Restart:   gen.SupervisorStrategyRestartTransient,
		},
	}
	return spec, nil
}

func (gr *Room) Init(p *ergo.Process, args ...interface{}) (state interface{}) {
	fmt.Printf("Initializing Go Room Process with room ID: %s\n", args[0])

	// Create a new GossipSub topic for this room
	topic, _ := ps.Join(args[0].(string))

	gr.topic = topic

	return RoomState{roomID: args[0].(string)}
}

func (gr *Room) HandleCast(message etf.Term, state interface{}) (string, interface{}) {
	fmt.Printf("Received message: %s\n", message)
	return "noreply", state
}

func (gr *Room) HandleCall(from ergo.ProcessFrom, message etf.Term, state interface{}) (string, etf.Term, interface{}) {
	fmt.Printf("Received message: %s from: %v\n", message, from)
	return "reply", etf.Term("ok"), state
}
