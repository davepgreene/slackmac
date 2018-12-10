package store

import (
	"errors"
	"fmt"
	"github.com/davepgreene/slackmac/propsd"
	log "github.com/sirupsen/logrus"
)

type PropsdStore struct {
	client *propsd.Client
	key string
}

func (p *PropsdStore) Get() string {
	resp, err := p.client.GetProperty(p.key)
	if err != nil {
		log.Error(err)
		return ""
	}

	return string(resp)
}

func NewPropsdStore(conf map[string]string) (Store, error) {
	key, ok := conf["key"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the propsd datastore", "key"))
	}

	return &PropsdStore{
		client: propsd.NewClient(propsd.DefaultEndpoint),
		key: key,
	}, nil
}
