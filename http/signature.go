package http

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"github.com/davepgreene/slackmac/errors"
	"github.com/davepgreene/slackmac/store"
	"github.com/davepgreene/slackmac/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

func signature(store store.Store, algorithm crypto.Hash) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		identifier := getCorrelationID(rw)
		metadata := map[string]interface{}{
			"identifier": identifier,
		}

		signature := r.Header.Get(slackSignatureHeader)

		log.WithFields(log.Fields{
			"identifier": identifier,
		}).Debugf("Request signature: %s", signature)

		requestTime, err := utils.EpochStringToTime(r.Header.Get(slackTimestampHeader))
		if err != nil {
			errors.ErrorWriter(errors.NewBadRequestError("Invalid date header", metadata), rw)
			return
		}

		// Save a copy of this request so we can put it back in the request once we consume it.
		buf, _ := ioutil.ReadAll(r.Body)
		body := bytes.NewBuffer(buf)
		deferredBody := ioutil.NopCloser(bytes.NewBuffer(buf))

		// Get the slack token and assemble the message we'll hash
		token := store.Get()
		message := slackSignatureVersion + ":" + strconv.FormatInt(requestTime.Unix(), 10) + ":" + body.String()
		log.Debugf("Message: %s", message)

		// Calculate the HMAC
		h := hmac.New(algorithm.New, []byte(token))
		h.Write([]byte(message))
		sha := hex.EncodeToString(h.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(slackSignatureVersion+ "=" + sha)) {
			errors.ErrorWriter(errors.NewAuthorizationError("Request HMAC is invalid", metadata), rw)
			return
		}

		r.Body = deferredBody
		next(rw, r)
	}
}
