package app

import (
	"github.com/bahner/go-myspace/config"
	"github.com/ergo-services/ergo/node"
)

var (
	log         = config.GetLogger()
	nodeCookie  = config.NodeCookie
	nodeName    = config.NodeName
	appName     = config.AppName
	version     = config.Version
	description = config.Description

	n node.Node
)
