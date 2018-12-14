package store

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"strings"
)

// Store represents an entity that can retrieve a slack token (or anything else) from an external source
type Store interface {
	Get() string
}

type factory func(conf map[string]string) (Store, error)

var storeFactories = make(map[string]factory)

// Register registers a Store type in the Store factory
func Register(name string, factory factory) {
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
	Register("propsd", newPropsdStore)
	Register("config", newConfigStore)
	Register("kms", newKMSStore)
	Register("secretsmanager", newSecretsManagerStore)
}

// CreateStore creates a Store
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
		return nil, fmt.Errorf("invalid Datastore name. Must be one of: %s", strings.Join(availableStores, ", "))
	}

	// Run the factory with the configuration.
	return storeFactory(conf)
}
