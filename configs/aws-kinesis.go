package configs

import (
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// KinesisConsumer ...
type KinesisConsumer struct {
	ResourceName string `json:"arn"`
	Params       struct {
		Group           string `json:"group"`
		EnableDebugLogs bool   `json:"enable_debug_logs"`
	} `json:"params"`
	resource resources.KinesisConsumer
	complete bool
}

// Bootstrap configuration.
func (cfg *KinesisConsumer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateKinesisConsumerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid kinesis consumer resource %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg KinesisConsumer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg KinesisConsumer) Resource() resources.KinesisConsumer {
	return cfg.resource
}

// KinesisProducer ...
type KinesisProducer struct {
	ResourceName string `json:"arn"`
	resource     resources.KinesisProducer
	complete     bool
}

// Bootstrap configuration.
func (cfg *KinesisProducer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateKinesisProducerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid kinesis producer resource; %w", err)
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg KinesisProducer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg KinesisProducer) Resource() resources.KinesisProducer {
	return cfg.resource
}
