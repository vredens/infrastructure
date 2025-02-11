package resources

import (
	"fmt"
)

// Kafka collection.
type Kafka struct {
	Clusters map[string]KafkaCluster `json:"clusters"`
}

// KafkaCluster data structure.
type KafkaCluster struct {
	Resource
	Brokers          []string          `json:"brokers"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	TopicPrefix      string            `json:"topic_prefix"`
	TopicSuffix      string            `json:"topic_suffix"`
	GroupPrefix      string            `json:"group_prefix"`
	GroupSuffix      string            `json:"group_suffix"`
	TopicTranslation map[string]string `json:"topic_translation"`
}

// Validate returns true if the resource is valid.
func (r KafkaCluster) Validate() error {
	if r.err != nil {
		return r.err
	}
	if len(r.Brokers) <= 0 {
		return fmt.Errorf("kafka consumer brokers can not be empty")
	}
	return r.err
}

// GroupNameFor a base group name will add any configured prefix or suffix to it.
func (r KafkaCluster) GroupNameFor(group string) string {
	return r.GroupPrefix + group + r.GroupSuffix
}

// TopicNameFor a base topic name will check for any topic translation or add any configured prefix or suffix to it.
func (r KafkaCluster) TopicNameFor(topic string) string {
	if translation, exists := r.TopicTranslation[topic]; exists {
		return translation
	}
	return r.TopicPrefix + topic + r.TopicSuffix
}

func (r KafkaCluster) sanitize() KafkaCluster {
	r.err = r.Validate()
	return r
}

// LocateKafkaClusterResource definition.
func (irl Locator) LocateKafkaClusterResource(arn string) KafkaCluster {
	var resource KafkaCluster
	var found bool
	var name, _, err = parse(arn, "messaging", "kafka", "clusters")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.Kafka.Clusters[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
