package proxy

import (
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"net/http"
	"net/url"
)

// Proxy represents a container for a oxy/forward.Forwarder
type Proxy struct {
	url *url.URL
	fwd *forward.Forwarder
}

// New creates a new Proxy that forwards to the provided url
func New(url string) *Proxy {
	fwd, _ := forward.New(forward.RoundTripper(&transport{
		RoundTripper: http.DefaultTransport,
		interceptorFunc: nil,
	}))

	return &Proxy{
		url: testutils.ParseURI(url),
		fwd: fwd,
	}
}

// ServeHTTP provides a consistent interface to register a Proxy as something that can handle HTTP requests
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// let us forward this request to another server
	r.URL = p.url
	p.fwd.ServeHTTP(rw, r)
}
