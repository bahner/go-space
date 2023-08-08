package topic

import (
	"myspace-pubsub/config"
	"sync"
)

var (
	log    = config.Log
	topics sync.Map
)
