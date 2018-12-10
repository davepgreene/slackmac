package store

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"strings"
)

type Store interface {
	Get() string
}

type Factory func(conf map[string]string) (Store, error)

var storeFactories = make(map[string]Factory)

func Register(name string, factory Factory) {
	if factory == nil {
		log.Fatalf("Datastore factory %s does not exist.", name)
	}
	_, registered := storeFactories[name]
	if registered {
		log.Errorf("Datastore factory %s already registered. Ignoring.", name)
	}
	storeFactories[name] = factory
}

func init() {
	Register("propsd", NewPropsdStore)
	Register("config", NewConfigStore)
	Register("kms", NewKMSStore)
	Register("secretsmanager", NewSecretsManagerStore)
}

func CreateStore(conf map[string]string) (Store, error) {
	// Get
	storeName, ok := conf["type"]
	if !ok {
		storeName = "config"
	}

	storeFactory, ok := storeFactories[storeName]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableStores := make([]string, len(storeFactories))
		for k := range storeFactories {
			availableStores = append(availableStores, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid Datastore name. Must be one of: %s", strings.Join(availableStores, ", ")))
	}

	// Run the factory with the configuration.
	return storeFactory(conf)
}
