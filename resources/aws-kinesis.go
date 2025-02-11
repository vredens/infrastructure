package resources

import "fmt"

// Kinesis resource collection.
type Kinesis struct {
	Consumers map[string]KinesisConsumer `json:"consumers"`
	Producers map[string]KinesisProducer `json:"producers"`
}

// KinesisConsumer resource data structure.
type KinesisConsumer struct {
	Resource
	AWS struct {
		Endpoint string `json:"endpoint"`
		Region   string `json:"region"`
	} `json:"aws"`
	Stream string `json:"stream"`
}

// Validate returns true if the resource is valid.
func (r KinesisConsumer) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Stream == "" {
		return fmt.Errorf("kinesis consumer stream can not be empty")
	}
	return r.err
}

func (r KinesisConsumer) sanitize() KinesisConsumer {
	r.err = r.Validate()
	return r
}

// LocateKinesisConsumerResource definition.
func (irl *Locator) LocateKinesisConsumerResource(arn string) KinesisConsumer {
	var resource KinesisConsumer
	var found bool
	var name, _, err = parse(arn, "messaging", "kinesis", "consumers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.Kinesis.Consumers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}

// KinesisProducer resource data structure.
type KinesisProducer struct {
	Resource
	AWS struct {
		Endpoint string `json:"endpoint"`
		Region   string `json:"region"`
	} `json:"aws"`
	Stream string `json:"stream"`
}

// Validate returns true if the resource is valid.
func (r KinesisProducer) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Stream == "" {
		return fmt.Errorf("kinesis producer stream can not be empty")
	}
	return r.err
}

func (r KinesisProducer) sanitize() KinesisProducer {
	r.err = r.Validate()
	return r
}

// LocateKinesisProducerResource definition.
func (irl *Locator) LocateKinesisProducerResource(arn string) KinesisProducer {
	var resource KinesisProducer
	var found bool
	var name, _, err = parse(arn, "messaging", "kinesis", "producers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.Kinesis.Producers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
