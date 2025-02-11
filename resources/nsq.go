package resources

import "fmt"

// NSQ resource configuration datastructure.
type NSQ struct {
	Consumers map[string]NSQConsumer `json:"consumers"`
	Producers map[string]NSQProducer `json:"producers"`
}

// NSQProducer resource configuration datastructure.
type NSQProducer struct {
	Resource
	NSQd        []string `json:"nsqd"`
	TopicPrefix string   `json:"topic_prefix"`
	TopicSuffix string   `json:"topic_suffix"`
}

func (r NSQProducer) Validate() error {
	if r.err != nil {
		return r.err
	}
	if len(r.NSQd) == 0 {
		return fmt.Errorf("no nsqd endpoints configured")
	}
	return r.err
}

// TopicNameFor a base topic name will check for any topic translation or add any configured prefix or suffix to it.
func (r NSQProducer) TopicNameFor(topic string) string {
	return r.TopicPrefix + topic + r.TopicSuffix
}

func (r NSQProducer) sanitize() NSQProducer {
	r.err = r.Validate()
	return r
}

// LocateNSQProducerResource ...
func (irl *Locator) LocateNSQProducerResource(arn string) NSQProducer {
	var resource NSQProducer
	var found bool
	var name, _, err = parse(arn, "messaging", "nsq", "producers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.NSQ.Producers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}

// NSQConsumer resource configuration datastructure.
type NSQConsumer struct {
	Resource
	NSQd          []string `json:"nsqd"`
	Lookupd       []string `json:"lookupd"`
	TopicPrefix   string   `json:"topic_prefix"`
	TopicSuffix   string   `json:"topic_suffix"`
	ChannelPrefix string   `json:"channel_prefix"`
	ChannelSuffix string   `json:"channel_suffix"`
}

func (r NSQConsumer) Validate() error {
	if r.err != nil {
		return r.err
	}
	if len(r.NSQd) == 0 && len(r.Lookupd) == 0 {
		return fmt.Errorf("no nsqd or lookupd endpoints configured")
	}
	return r.err
}

// ChannelNameFor a base channel name will add any configured prefix or suffix to it.
func (resource NSQConsumer) ChannelNameFor(channel string) string {
	return resource.ChannelPrefix + channel + resource.ChannelSuffix
}

// TopicNameFor a base topic name will check for any topic translation or add any configured prefix or suffix to it.
func (resource NSQConsumer) TopicNameFor(topic string) string {
	return resource.TopicPrefix + topic + resource.TopicSuffix
}

func (resource NSQConsumer) sanitize() NSQConsumer {
	resource.err = resource.Validate()
	return resource
}

// LocateNSQConsumerResource ...
func (irl *Locator) LocateNSQConsumerResource(arn string) NSQConsumer {
	var resource NSQConsumer
	var found bool
	var name, _, err = parse(arn, "messaging", "nsq", "consumers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.NSQ.Consumers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
