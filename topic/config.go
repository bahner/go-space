package topic

import (
	"sync"

	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/global"
)

var (
	topics sync.Map

	log = config.Log
	ps  = *global.PubSubService
)
