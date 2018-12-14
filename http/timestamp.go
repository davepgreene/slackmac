package http

import (
	"fmt"
	"github.com/davepgreene/slackmac/errors"
	"github.com/davepgreene/slackmac/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var requiredHeaders = [2]string{slackTimestampHeader, slackSignatureHeader}

func timestamp(skew time.Duration) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		metadata := map[string]interface{} {
			"identifier": getCorrelationID(rw),
		}
		fields := log.Fields{
			"identifier": metadata["identifier"],
		}

		for _, header := range requiredHeaders {
			if r.Header.Get(header) == "" {
				errors.ErrorWriter(errors.NewBadRequestError(fmt.Sprintf("Missing header: %s", header), metadata), rw)
				return
			}
		}

		requestTime, err := utils.EpochStringToTime(r.Header.Get(slackTimestampHeader))
		if err != nil {
			errors.ErrorWriter(errors.NewBadRequestError("Invalid date header", metadata), rw)
			return
		}

		// Verify that the request date is close to $NOW
		timeSinceRequest := time.Since(requestTime)

		log.WithFields(fields).Debugf("Request Time: %s", requestTime)
		log.WithFields(fields).Debug("Time since request: %s", timeSinceRequest)
		log.WithFields(fields).Debugf("Timestamp skew %s", timeSinceRequest)

		if timeSinceRequest > skew {
			errors.ErrorWriter(errors.NewAuthorizationError("Request has expired", metadata), rw)
			return
		}

		next(rw, r)
	}
}
