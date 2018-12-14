package http

import (
	prox "github.com/davepgreene/slackmac/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/stream"
	"net/http"
)

func proxy(url string) http.Handler {
	proxy := prox.New(url)

	s, err := stream.New(proxy)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
