package config

import (
	"crypto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MetricsClientConfig struct {
	Host string
	Port int
	Prefix string
	Enabled bool
}

// Defaults generates a set of default configuration options
func Defaults() {
	viper.SetDefault("listen", map[string]interface{}{
		"bind": "0.0.0.0",
		"port": 9300,
		"limit": "10mb",
	})

	viper.SetDefault("log", map[string]interface{}{
		"level": log.InfoLevel,
		"json": true,
		"requests": true,
	})

	viper.SetDefault("correlation", map[string]interface{}{
		"enabled": true,
		"header": "X-Request-Identifier",
	})

	viper.SetDefault("slack", map[string]interface{}{
		"algorithm": "SHA256",
		"skew": "5m",
	})

	viper.SetDefault("store", map[string]interface{} {
		"type": "config",
		"key": nil,
	})

	viper.SetDefault("service", map[string]interface{}{
		"port":     9301,
		"hostname": "127.0.0.1",
		"limit":    "10mb",
		"protocol": "http://",
	})

	viper.SetDefault("metrics", map[string]interface{}{
		"enabled": true,
	})
	viper.SetDefault("metrics.client", map[string]interface{}{
		"host":   "localhost",
		"port":   8125,
		"prefix": "slackmac.",
	})
}

// Metrics fixes an issue where config files make nested values in the same
// map disappear.
func Metrics() MetricsClientConfig {
	var conf MetricsClientConfig
	enabled := viper.GetBool("metrics.enabled")

	if !enabled {
		conf.Enabled = false
		return conf
	}

	err := viper.UnmarshalKey("metrics.client", &conf)
	if err != nil {
		log.Error("Unable to unmarshal metrics configuration. No metrics will be captured")
		conf.Enabled = false
	} else {
		conf.Enabled = true
	}
	return conf
}

var SUPPORTED_ALGORITHMS = map[string]interface{}{
	"SHA256": crypto.SHA256,
}

var SUPPORT_ALGORITHMS_LOOKUP = map[crypto.Hash]string {
	crypto.SHA256: "SHA256",
}

