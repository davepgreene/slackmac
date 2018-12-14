package store

import (
	"fmt"
	"github.com/spf13/viper"
)

type configStore struct {
	key string
}

// Get retrieves data from the local config
func (c *configStore) Get() string {
	return viper.GetString(c.key)
}

func newConfigStore(conf map[string]string) (Store, error) {
	key, ok := conf["key"]
	if !ok {
		return nil, fmt.Errorf("%s is required for the config datastore", "key")
	}

	return &configStore{
		key: key,
	}, nil
}

