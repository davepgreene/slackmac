package store

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type ConfigStore struct {
	key string
}
func (c *ConfigStore) Get() string {
	return viper.GetString(c.key)
}

func NewConfigStore(conf map[string]string) (Store, error) {
	key, ok := conf["key"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the config datastore", "key"))
	}

	return &ConfigStore{
		key: key,
	}, nil
}

