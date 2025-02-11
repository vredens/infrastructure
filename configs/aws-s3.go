package configs

import "github.com/vredens/infrastructure/resources"

// S3Manager defines the configuration required for an s3 client.
type S3Manager struct {
	ResourceName string `json:"arn"`
	resource     resources.S3Manager
	complete     bool
}

// Bootstrap configuration.
func (cfg *S3Manager) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateS3ManagerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return err
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg S3Manager) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg S3Manager) Resource() resources.S3Manager {
	return cfg.resource
}
