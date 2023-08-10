package topic

import (
	"sync"

	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/global"
)

var (
	log             = config.Log
	topics          sync.Map
	myspaceNodeName = config.MyspaceNodeName
	ps              = *global.PubSubService
)
