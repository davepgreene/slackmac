package http

import (
	"github.com/davepgreene/slackmac/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func metricsMiddleware (rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	instance, err := utils.Metrics()
	// Move on to the next middleware if the metrics client is disabled
	if err != nil {
		next(rw, r)
	}

	err = instance.Incr("request.incoming", nil, 1)
	if err != nil {
		log.Error(err)
	}

	err = instance.Gauge("request.active", 1, nil, 1)
	if err != nil {
		log.Error(err)
	}

	next(rw, r)
}
