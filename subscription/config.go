package subscription

import (
	"github.com/bahner/go-myspace/config"
	"github.com/bahner/go-myspace/global"
)

var (
	log             = config.Log
	myspaceNodeName = config.MyspaceNodeName
	ps              = *global.PubSubService
)
