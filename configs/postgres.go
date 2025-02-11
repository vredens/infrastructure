package configs

import (
	"errors"
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// Postgres for a new postgres SQL connection.
type Postgres struct {
	ResourceName string `json:"arn"`
	Params       struct {
		// MaxIdleConns is the maximum number of idle connections kept in the connection pool.
		MaxIdleConns int `json:"max_idle_conns"`
		// MaxOpenConns is the maximum number of open connections to the database.
		MaxOpenConns int `json:"max_open_conns"`
	} `json:"params"`
	resource resources.Postgres
	complete bool
}

// Bootstrap configuration.
func (cfg *Postgres) Bootstrap(provider resources.Provider) error {
	if cfg.Params.MaxIdleConns <= 0 {
		cfg.Params.MaxIdleConns = 2
	}
	if cfg.Params.MaxOpenConns <= 0 {
		cfg.Params.MaxOpenConns = 10
	}

	cfg.resource = provider.Locator().LocatePostgresResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("could not locate infrastructure resource for %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Valid returns an error if the configuration is NOT valid.
func (cfg Postgres) Valid() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource returns the infrastructure resource located by a previous call to Complete.
func (cfg Postgres) Resource() resources.Postgres {
	return cfg.resource
}

// GetDSN for this configuration. Returns an empty string if configuration is incomplete.
func (cfg Postgres) GetDSN() string {
	if !cfg.complete {
		return ""
	}
	return cfg.resource.GetDSN()
}

// GetFullDSN which includes extra connection parameters in the format <dsn>?<key>=<value>. Returns an empty string if configuration is incomplete.
func (cfg Postgres) GetFullDSN() string {
	if !cfg.complete {
		return ""
	}
	return cfg.resource.GetFullDSN()
}

// PostgresListenerConfig for a postgres pg_notify listener.
type PostgresListenerConfig struct {
	ResourceName string                       `json:"arn"`
	Params       PostgresListenerConfigParams `json:"params"`
	resource     resources.Postgres
	complete     bool
}

// PostgresListenerConfigParams for a postgres listener resource.
type PostgresListenerConfigParams struct {
	Channel string `json:"channel"`
}

// Bootstrap configuration.
func (cfg *PostgresListenerConfig) Bootstrap(provider resources.Provider) error {
	if cfg.Params.Channel == "" {
		return errors.New("invalid channel")
	}

	cfg.resource = provider.Locator().LocatePostgresResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("could not locate infrastructure resource for %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Valid returns an error if the configuration is NOT valid.
func (cfg PostgresListenerConfig) Valid() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource returns the infrastructure resource located by a previous call to Complete.
func (cfg PostgresListenerConfig) Resource() resources.Postgres {
	return cfg.resource
}
