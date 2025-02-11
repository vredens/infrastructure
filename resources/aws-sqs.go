package resources

import "fmt"

// SQSResource collection.
type SQSResource struct {
	Consumers map[string]SQSConsumerResource `json:"consumers"`
	Producers map[string]SQSProducerResource `json:"producers"`
}

// SQSConsumerResource data structure.
type SQSConsumerResource struct {
	Resource
	AWS struct {
		Endpoint string `json:"endpoint"`
		Region   string `json:"region"`
	} `json:"aws"`
	Queue string `json:"queue"`
}

// Validate returns true if the resource is valid.
func (r SQSConsumerResource) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Queue == "" {
		return fmt.Errorf("sqs consumer queue can not be empty")
	}
	return r.err
}

func (r SQSConsumerResource) sanitize() SQSConsumerResource {
	r.err = r.Validate()
	return r
}

// LocateSQSConsumerResource definition.
func (irl *Locator) LocateSQSConsumerResource(arn string) SQSConsumerResource {
	var resource SQSConsumerResource
	var found bool
	var name, _, err = parse(arn, "messaging", "sqs", "consumers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.SQS.Consumers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}

// SQSProducerResource data structure.
type SQSProducerResource struct {
	Resource
	AWS struct {
		Endpoint string `json:"endpoint"`
		Region   string `json:"region"`
	} `json:"aws"`
	Queue string `json:"queue"`
}

// Validate returns true if the resource is valid.
func (r SQSProducerResource) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Queue == "" {
		return fmt.Errorf("sqs producer queue can not be empty")
	}
	return r.err
}

func (r SQSProducerResource) sanitize() SQSProducerResource {
	r.err = r.Validate()
	return r
}

// LocateSQSProducerResource definition.
func (irl *Locator) LocateSQSProducerResource(arn string) SQSProducerResource {
	var resource SQSProducerResource
	var found bool
	var name, _, err = parse(arn, "messaging", "sqs", "producers")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Messaging.SQS.Producers[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
