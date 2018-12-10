package proxy

import (
	"bytes"
	"github.com/davepgreene/slackmac/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type transport struct {
	http.RoundTripper
	interceptorFunc func(body []byte) []byte
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timer := time.Now()
	// Before request
	resp, err = t.RoundTripper.RoundTrip(req)
	// After response back

	if err != nil {
		return nil, err
	}


	responseDuration := time.Since(timer)
	instance, err := utils.Metrics()
	if err == nil {
		err = instance.Timing("proxy.duration", responseDuration, nil, 1)
		if err != nil {
			log.Error(err)
		}
	}

	// We might not need to read the response body at all
	if t.interceptorFunc != nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		log.Debug("Executed transport interceptor func.")
		b = t.interceptorFunc(b)

		body := ioutil.NopCloser(bytes.NewReader(b))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	}

	return resp, nil
}
