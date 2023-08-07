package main

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
)

type Myspace struct {
	gen.Server
}

func createMyspace() gen.ServerBehavior {

	return &Myspace{}

}

func (gr *Myspace) Init(sp *gen.ServerProcess, args ...etf.Term) error {

	log.Info("Initializing Myspace Dispatcher Process")

	return nil

}

func (gr *Myspace) HandleCast(server_procces *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Infof("Creating new room with no reply: %s\n", message)
	return gen.ServerStatusOK
}

func (gr *Myspace) HandleCall(serverProcess *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	log.Infof("Creating new room with a reply: %s\n", message)
	msg := fmt.Sprintf("Creating new room with a reply: %s\n", message)
	go spawnAndRegisterRoom(message.(string))
	return msg, gen.ServerStatusOK
}

func (gr *Myspace) HandleInfo(serverProcess *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("Received message: %s\n", message)
	return gen.ServerStatusOK
}
