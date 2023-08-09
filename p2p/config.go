package p2p

import (
	"myspace-pubsub/config"

	logging "github.com/ipfs/go-log"
)

var (
	rendezvous = config.Rendezvous
	log        = config.Log
	name       = config.AppName
	loglevel   = config.LogLevel
)

func initLogging() {
	logging.SetLogLevel(name, *loglevel)
}
