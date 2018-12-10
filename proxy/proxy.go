package proxy

import (
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"net/http"
	"net/url"
)

type Proxy struct {
	url *url.URL
	fwd *forward.Forwarder
}

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

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// let us forward this request to another server
	r.URL = p.url
	p.fwd.ServeHTTP(rw, r)
}
