package node

import (
	"flag"

	"go.deanishe.net/env"
)

var (
	defaultNodeCookie = env.Get("MYSPACE_NODE_COOKIE", "myspace")
	defaultNodeName   = env.Get("MYSPACE_NODE_NAME", "pubsub@localhost")
)

var (
	nodeCookie = flag.String("nodecookie", defaultNodeCookie, "Secret shared by all erlang nodes in the cluster")
	nodeName   = flag.String("nodename", defaultNodeName, "Name of the erlang node")
)
