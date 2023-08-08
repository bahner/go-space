package app

import (
	"myspace-pubsub/config"

	"github.com/ergo-services/ergo/node"
)

var (
	log         = config.Log
	nodeCookie  = config.NodeCookie
	nodeName    = config.NodeName
	appName     = config.AppName
	version     = config.Version
	description = config.Description
	n           node.Node
)
