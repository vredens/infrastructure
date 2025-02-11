package configs

import (
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// SQSConsumer ...
type SQSConsumer struct {
	ResourceName string `json:"arn"`
	Params       struct {
		Workers              int  `json:"workers"`
		MaxMessagesPerWorker int  `json:"max_messages_per_worker"`
		RawMessage           bool `json:"raw_message"`
		EnableDebugLogs      bool `json:"enable_debug_logs"`
	} `json:"params"`
	resource resources.SQSConsumerResource
	complete bool
}

// Bootstrap configuration.
func (cfg *SQSConsumer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateSQSConsumerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid sqs consumer resource %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg SQSConsumer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg SQSConsumer) Resource() resources.SQSConsumerResource {
	return cfg.resource
}

// SQSProducer ...
type SQSProducer struct {
	ResourceName string `json:"arn"`
	resource     resources.SQSProducerResource
	complete     bool
}

// Bootstrap configuration.
func (cfg *SQSProducer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateSQSProducerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid sqs producer resource %s; %w", cfg.ResourceName, err)
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg SQSProducer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg SQSProducer) Resource() resources.SQSProducerResource {
	return cfg.resource
}
