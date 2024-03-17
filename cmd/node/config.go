package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defautSpaceNodeName      = "space@localhost"
	defaultNodeCookie        = "spacecookie"
	defaultNodeName          = "pubsub@localhost"
	defaultNodeDebugInterval = time.Second * 60
)

func init() {

	// Erlang node config
	pflag.String("spacenode", defautSpaceNodeName, "Name of the node running the actual SPACE")
	viper.BindPFlag("node.space", pflag.Lookup("spacenode"))
	viper.SetDefault("node.space", defautSpaceNodeName)

	pflag.String("nodecookie", defaultNodeCookie, "Secret shared between erlang nodes in the cluster")
	viper.BindPFlag("node.cookie", pflag.Lookup("nodecookie"))
	viper.SetDefault("node.cookie", defaultNodeCookie)

	pflag.String("nodename", defaultNodeName, "Name of the erlang node")
	viper.BindPFlag("node.name", pflag.Lookup("nodename"))
	viper.SetDefault("node.name", defaultNodeName)

	pflag.Duration("node_debug_interval", defaultNodeDebugInterval, "Interval for debug output")
	viper.BindPFlag("node.debug_interval", pflag.Lookup("_node_debug_interval"))
	viper.SetDefault("node.debug_interval", defaultNodeDebugInterval)

}
