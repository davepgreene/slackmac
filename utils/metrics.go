package utils

import (
	"errors"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/davepgreene/slackmac/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

var instance *statsd.Client
var err error
var once sync.Once

// Metrics creates a singleton client for submitting metrics to DogStatsD
func Metrics() (*statsd.Client, error) {
	once.Do(func() {
		metricsClientConf := config.Metrics()
		if !metricsClientConf.Enabled {
			err = errors.New("metrics collection has been disabled")
			log.Error(err)
			instance = nil
		} else {
			conn := fmt.Sprintf("%s:%d", metricsClientConf.Host, metricsClientConf.Port)
			instance, err = statsd.New(conn)
			if err != nil {
				log.Info("Unable to initialize statsd client.")
				log.Fatal(err)
			}
			instance.Namespace = metricsClientConf.Prefix
		}
	})
	return instance, err
}
