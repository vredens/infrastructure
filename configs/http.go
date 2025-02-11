package configs

import "github.com/vredens/infrastructure/resources"

// Webservice configuration for connecting to a webservice.
type Webservice struct {
	ResourceName string `json:"arn"`
	Params       struct {
		Connection HTTPConnection `json:"connection"`
		// Timeout for new connections, in milliseconds.
		Timeout int `json:"timeout"`
	} `json:"params"`
	resource resources.Webservice
	complete bool
}

// HTTPConnection for finetuning the connection.
type HTTPConnection struct {
	// MaxIdle is the maximum number of idle connections.
	MaxIdle int `json:"max_idle"`
	// MaxPerHost is the maximum number of connections per host.
	MaxPerHost int `json:"max_per_host"`
	// KeepAlive for connections, in seconds.
	KeepAlive int `json:"keep_alive"`
}

// Bootstrap configuration.
func (cfg *Webservice) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateWebserviceResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return err
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg Webservice) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg Webservice) Resource() resources.Webservice {
	return cfg.resource
}
