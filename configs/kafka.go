package configs

import (
	"fmt"

	"github.com/vredens/infrastructure/resources"
)

// KafkaCluster configuration.
type KafkaCluster struct {
	ResourceName string `json:"arn"`
	resource     resources.KafkaCluster
	complete     bool
}

// Bootstrap configuration.
func (cfg *KafkaCluster) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateKafkaClusterResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid kafka cluster resource; %w", err)
	}

	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg KafkaCluster) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg KafkaCluster) Resource() resources.KafkaCluster {
	return cfg.resource
}

// TopicNameFor returns the topic name to use by applying any configured prefix and/or suffix.
// This allows to easily setup topic names per environment such as adding an environment suffix.
// Which is useful when using a single Kafka cluster in non-production environments.
func (cfg KafkaCluster) TopicNameFor(topic string) string {
	return cfg.resource.TopicNameFor(topic)
}

// GroupNameFor returns the consumer group name to use by applying any configured prefix and/or suffix.
// This allows to easily setup consumer group names per environment such as adding an environment suffix.
// Which is useful when using a single Kafka cluster in non-production environments.
func (cfg KafkaCluster) GroupNameFor(group string) string {
	return cfg.resource.GroupPrefix + group + cfg.resource.GroupSuffix
}

// KafkaConsumer ...
type KafkaConsumer struct {
	ResourceName string `json:"arn"`
	// Topic from which to consume messages.
	Topic string `json:"topic"`
	// Group is the consumer group name.
	Group string `json:"group"`
	// InitialOffset can be relative (-100, +200) or timestamp (@2000-01-02T03:04:05.006Z).
	// This parameter is dependent on driver support.
	InitialOffset string `json:"initial_offset"`
	resource      resources.KafkaCluster
	complete      bool
}

// Bootstrap configuration.
func (cfg *KafkaConsumer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateKafkaClusterResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid kafka consumer resource %s; %w", cfg.ResourceName, err)
	}
	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg KafkaConsumer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg KafkaConsumer) Resource() resources.KafkaCluster {
	return cfg.resource
}

// TopicName for this configuration which includes any prefix/suffix specified in the infra resource.
func (cfg KafkaConsumer) TopicName() string {
	return cfg.resource.TopicNameFor(cfg.Topic)
}

// GroupName for this configuration which includes any prefix/suffix specified in the infra resource.
func (cfg KafkaConsumer) GroupName() string {
	return cfg.GroupNameFor(cfg.Group)
}

// TopicNameFor returns the topic name to use by applying any configured prefix and/or suffix.
// This allows to easily setup topic names per environment such as adding an environment suffix.
// Which is useful when using a single Kafka cluster in non-production environments.
func (cfg KafkaConsumer) TopicNameFor(topic string) string {
	return cfg.resource.TopicNameFor(topic)
}

// GroupNameFor returns the consumer group name to use by applying any configured prefix and/or suffix.
// This allows to easily setup consumer group names per environment such as adding an environment suffix.
// Which is useful when using a single Kafka cluster in non-production environments.
func (cfg KafkaConsumer) GroupNameFor(group string) string {
	return cfg.resource.GroupPrefix + group + cfg.resource.GroupSuffix
}

// KafkaProducer ...
type KafkaProducer struct {
	ResourceName string `json:"arn"`
	Topic        string `json:"topic"`
	resource     resources.KafkaCluster
	complete     bool
}

// Bootstrap configuration.
func (cfg *KafkaProducer) Bootstrap(provider resources.Provider) error {
	cfg.resource = provider.Locator().LocateKafkaClusterResource(cfg.ResourceName)
	if err := cfg.resource.Validate(); err != nil {
		return fmt.Errorf("invalid kafka producer resource; %w", err)
	}
	cfg.complete = true

	return nil
}

// Validate returns an error if the configuration is NOT valid.
func (cfg KafkaProducer) Validate() error {
	if !cfg.complete {
		return ErrConfigNotBootstrapped
	}
	return cfg.resource.Validate()
}

// Resource with the infrastructure configuration.
func (cfg KafkaProducer) Resource() resources.KafkaCluster {
	return cfg.resource
}

// TopicName for this configuration which includes any prefix/suffix specified in the infra resource.
func (cfg KafkaProducer) TopicName() string {
	return cfg.resource.TopicNameFor(cfg.Topic)
}

// TopicName returns the topic name to use by applying any configured prefix and/or suffix.
// This allows to easily setup topic names per environment such as adding an environment suffix.
// Which is useful when using a single Kafka cluster in non-production environments.
func (cfg KafkaProducer) TopicNameFor(topic string) string {
	return cfg.resource.TopicNameFor(topic)
}
