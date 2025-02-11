package configs

import (
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// Elasticsearch allows changing some of the default connector options.
type Elasticsearch struct {
	ResourceName string `json:"arn"`
	Params       struct {
		// Timeout for HTTP connections. Must be between 1 and 300 seconds. Defaults to 5 seconds.
		Timeout            int `json:"timeout"`
		MaxIdleConnections int `json:"max_idle_connections"`
	} `json:"params"`
	complete bool
	resource resources.Elasticsearch
}

// Bootstrap the configuration settings by finding the infrastructure resource and sanitizing the configuration parameters.
func (cfg *Elasticsearch) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateElasticResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("failed to locate infrastructure resource for %s; %w", cfg.ResourceName, err)
	}

	if cfg.Params.Timeout < 1 || cfg.Params.Timeout >= 300 {
		cfg.Params.Timeout = 5
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg Elasticsearch) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource for this configuration. Requires previous call to Bootstrap.
func (cfg Elasticsearch) Resource() resources.Elasticsearch {
	return cfg.resource
}
