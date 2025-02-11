package configs

import (
	"errors"
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// AlgoliaConfig data structure for algolia connections.
type AlgoliaConfig struct {
	ResourceName string `json:"arn"`
	complete     bool
	resource     resources.Algolia
}

// Bootstrap the configuration settings by finding the infrastructure resource and sanitizing the configuration parameters.
func (cfg *AlgoliaConfig) Bootstrap(provider resources.Provider) error {
	if cfg.complete {
		return errors.New("configuration already bootstrapped")
	}
	cfg.resource = provider.Locator().LocateAlgoliaResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid infrastructure resource for %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Valid returns an error if the configuration is NOT valid.
func (cfg AlgoliaConfig) Valid() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource for this configuration. Requires previous call to Bootstrap.
func (cfg AlgoliaConfig) Resource() resources.Algolia {
	return cfg.resource
}

// IndexNameFor computes the final index name for the provided index using the configured index prefix.
// If configuration is not yet bootstrapped an empty string is returned instead.
func (cfg AlgoliaConfig) IndexNameFor(index string) string {
	if !cfg.complete {
		return ""
	}
	return cfg.resource.IndexPrefix + index
}
