package configs

import "github.com/vredens/infrastructure/resources"

// SFTP configuration required for an sftp client.
type SFTP struct {
	ResourceName string `json:"arn"`
	resource     resources.SFTP
	complete     bool
}

// Bootstrap configuration.
func (cfg *SFTP) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateSFTPResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return err
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg SFTP) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg SFTP) Resource() resources.SFTP {
	return cfg.resource
}
