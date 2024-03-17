package main

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/spf13/viper"
)

func configTemplate(identity string, node string) map[string]interface{} {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	return map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"location": config.ActorLocation(),
			"nick":     config.ActorNick(),
		},
		"db": map[string]interface{}{
			"file": config.DefaultDbFile,
		},
		// Use default log settings, so as not to pick up debug log settings
		"log": map[string]interface{}{
			"level": viper.GetString("log.level"),
			"file":  viper.GetString("log.file"),
		},
		// NB! This is a cross over from go-ma
		"api": map[string]interface{}{
			// This must be set corretly for generation to work
			"maddr": viper.GetString("api.maddr"),
		},
		"http": map[string]interface{}{
			"socket": config.HttpSocket(),
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     config.P2PPort(),
			"connmgr": map[string]interface{}{
				"low-watermark":  config.P2PConnmgrLowWatermark(),
				"high-watermark": config.P2PConnmgrHighWatermark(),
				"grace-period":   config.P2PConnMgrGracePeriod(),
			},
			"discovery": map[string]interface{}{
				"advertise-ttl":      config.P2PDiscoveryAdvertiseTTL(),
				"advertise-limit":    config.P2PDiscoveryAdvertiseLimit(),
				"advertise-interval": config.P2PDiscoveryAdvertiseInterval(),
				"dht":                config.P2PDiscoveryDHT(),
				"mdns":               config.P2PDiscoveryMDNS(),
			},
		},
		"node": map[string]interface{}{
			"cookie":         defaultNodeCookie,
			"name":           defaultNodeName,
			"debug_interval": defaultNodeDebugInterval,
			"space":          defautSpaceNodeName,
		},
	}
}
