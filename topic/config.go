package topic

import (
	"sync"

	"github.com/bahner/go-myspace/config"
)

var (
	log             = config.Log
	topics          sync.Map
	myspaceNodeName = config.MyspaceNodeName
)
