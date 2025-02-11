package configs

import "github.com/vredens/infrastructure/resources"

// Dynamo defines the configuration required for a dynamo client.
type Dynamo struct {
	ResourceName string `json:"arn"`
	resource     resources.Dynamo
	complete     bool
}

// Bootstrap configuration.
func (cfg *Dynamo) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateDynamoResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return err
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg Dynamo) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg Dynamo) Resource() resources.Dynamo {
	return cfg.resource
}
