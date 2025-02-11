package configs

import "github.com/vredens/infrastructure/resources"

// Redis configuration is used to convert a redis resource location into a RedisResource.
type Redis struct {
	ResourceName string `json:"arn"`
	resource     resources.Redis
	complete     bool
}

// Bootstrap configuration.
func (cfg *Redis) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateRedisResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return err
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg Redis) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg Redis) Resource() resources.Redis {
	return cfg.resource
}
