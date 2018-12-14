package store

import (
	"fmt"
	"github.com/davepgreene/slackmac/propsd"
	log "github.com/sirupsen/logrus"
	"time"
)

type propsdStore struct {
	client *propsd.Client
	key string
}

// Get retrieves data from Propsd
func (p *propsdStore) Get() string {
	resp, err := p.client.GetProperty(p.key)
	if err != nil {
		log.Error(err)
		return ""
	}

	return string(resp)
}

func newPropsdStore(conf map[string]string) (Store, error) {
	key, ok := conf["key"]
	if !ok {
		return nil, fmt.Errorf("%s is required for the propsd datastore", "key")
	}
	var endpoint string
	endpoint, ok = conf["endpoint"]
	if !ok {
		endpoint = propsd.DefaultEndpoint
	}

	var timeout = propsd.DefaultTimeout
	if timeoutStr, ok := conf["timeout"]; ok {
		t, err := time.ParseDuration(timeoutStr)
		if err == nil {
			timeout = t
		}
	}

	return &propsdStore{
		client: propsd.NewClient(endpoint, timeout),
		key: key,
	}, nil
}
