package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// GetLogLevel retrieves the desired log level from settings.
//
// NOTE: This should only be called after viper initializes
func GetLogLevel() log.Level {
	if lvl, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		return lvl
	}

	log.Info("Unable to parse log level in settings. Defaulting to INFO.")
	return log.InfoLevel
}

// GetLogFormatter  retrieves the desired log formatter from settings.
//
// NOTE: This should only be called after viper initializes
func GetLogFormatter() log.Formatter {
	fmt := viper.GetBool("log.json")
	if fmt {
		return &log.JSONFormatter{}
	}
	return &prefixed.TextFormatter{}
}
