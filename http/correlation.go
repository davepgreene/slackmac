package http

import (
	"github.com/davepgreene/slackmac/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func GetCorrelationId(rw http.ResponseWriter) string {
	if viper.GetBool("correlation.enabled") {
		return rw.Header().Get(viper.GetString("correlation.header"))
	}

	return ""
}

func correlationMiddleware(correlation map[string]interface{}) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	header, ok := correlation["header"].(string)
	if ok == false {
		log.Fatal("Missing required parameter `header`!")
	}

	log.Infof("Using Correlation-Identifier %s", header)

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id := r.Header.Get(header)
		// If no header is set, create a correlation id
		if id == "" {
			id = utils.CreateCorrelationID()
			rw.Header().Add(header, id)

			// Add the correlation ID to the request so that downstream consumers can track it
			r.Header.Set(header, id)

			log.WithFields(log.Fields{
				"identifier": id,
			}).Debugf("Setting Correlation-Identifier header %s:%s", header, id)

			next(rw, r)
		} else {
			log.WithFields(log.Fields{
				"identifier": id,
			}).Debugf("Found Correlation-Identifier header %s:%s", header, id)

			next(rw, r)
		}
	}
}
