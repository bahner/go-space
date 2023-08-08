package topic

import (
	"myspace-pubsub/config"
	"myspace-pubsub/p2p"
	"sync"
)

var (
	log    = config.Log
	ps     = p2p.PubSubService
	topics sync.Map
)
