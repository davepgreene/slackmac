package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
	"strings"
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
		return &UppercaseJSONFormatter{}
	}
	return &prefixed.TextFormatter{}
}

// UppercaseJSONFormatter takes JSON formatted logs and uppercases the log level
type UppercaseJSONFormatter struct {
	log.JSONFormatter
}

// Format renders a single log entry
func (f *UppercaseJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	d := make(map[string]interface{})
	b, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	if val, ok := d["level"]; ok {
		d["level"] = strings.ToUpper(val.(string))
	}

	b, err = json.Marshal(d)
	return append(b, "\n"...), err
}
