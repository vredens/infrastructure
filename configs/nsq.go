package configs

import (
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// NSQProducer data structure.
type NSQProducer struct {
	ResourceName string `json:"arn"`
	// Topic for sending messages. Should be in the format <system>.<schema version>.<topic>[.<partition|priority>]
	Topic    string `json:"topic"`
	resource resources.NSQProducer
	complete bool
}

// Bootstrap configuration.
func (cfg *NSQProducer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateNSQProducerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid nsq producer resource; %w", err)
	}
	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg NSQProducer) Validate() error {
	if !cfg.complete {
		return fmt.Errorf("configuration not bootstrapped")
	}
	return nil
}

// Resource with the infrastructure configuration.
func (cfg NSQProducer) Resource() resources.NSQProducer {
	return cfg.resource
}

// TopicName returns the topic name to use by applying any configured prefix and/or suffix.
func (cfg NSQProducer) TopicName() string {
	return cfg.resource.TopicNameFor(cfg.Topic)
}

// NSQConsumer data structure.
type NSQConsumer struct {
	ResourceName string `json:"arn"`
	// Topic is prefixed with '<env>.' if topic prefix is enabled.
	Topic string `json:"topic" validate:"nonzero"`
	// Channel is prefixed with '<env>.<system>.' if channel prefix is enabled.
	Channel         string `json:"channel" validate:"nonzero"`
	Concurrency     int    `json:"concurrency"`
	RequeueDelay    int    `json:"requeue_delay"`
	MaxRequeueDelay int    `json:"max_requeue_delay"`
	MaxAttempts     int    `json:"max_attempts"`
	resource        resources.NSQConsumer
	complete        bool
}

// Bootstrap configuration.
func (cfg *NSQConsumer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateNSQConsumerResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid nsq consumer resource %s; %w", cfg.ResourceName, err)
	}
	if cfg.Topic == "" {
		return fmt.Errorf("no topic configured")
	}
	if cfg.Channel == "" {
		return fmt.Errorf("no channel configured")
	}
	if channel := cfg.resource.ChannelNameFor(cfg.Channel); len(channel) > 64 {
		return fmt.Errorf("channel name must be under 64 characters total [channel:%s][size:%d]", channel, len(channel))
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg NSQConsumer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg NSQConsumer) Resource() resources.NSQConsumer {
	return cfg.resource
}

// TopicName returns the topic name to use by applying any configured prefix and/or suffix.
func (cfg NSQConsumer) TopicName() string {
	return cfg.resource.TopicNameFor(cfg.Topic)
}

// ChannelName returns the channel name to use by applying any configured prefix and/or suffix.
func (cfg NSQConsumer) ChannelName() string {
	return cfg.resource.ChannelNameFor(cfg.Channel)
}
