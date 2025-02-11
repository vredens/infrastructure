package resources

import (
	"fmt"
	"strings"
)

const (
	resourcePrefix = "arn://"
)

var (
	ErrConfigNotBootstrapped = fmt.Errorf("configuration not bootstrapped")
	ErrResourceNotFound      = fmt.Errorf("resource not found")
)

type Provider interface {
	// Locator gives access to the infrastructure configuration for implementing your own providers.
	Locator() *Locator
	// SystemName your process is running in.
	SystemName() string
	// ComponentName your process is running in.
	ComponentName() string
	// Environment your process is running in.
	Environment() string
	// RenderSecret ...
	RenderSecret(value string) string
}

// Locator of resources.
type Locator struct {
	Cloud       Cloud                 `json:"cloud"`
	Databases   Databases             `json:"storage"`
	Messaging   Messaging             `json:"messaging"`
	Webservices map[string]Webservice `json:"webservices"`
	provider    Provider
}

func (loc *Locator) SetProvider(provider Provider) {
	loc.provider = provider
}

type Cloud struct {
	AWS map[string]AWSSession `json:"aws"`
}

// Databases configuration datastructure.
type Databases struct {
	Algolia       map[string]Algolia       `json:"algolia"`
	Elasticsearch map[string]Elasticsearch `json:"elasticsearch"`
	Postgres      map[string]Postgres      `json:"postgres"`
	Redis         map[string]Redis         `json:"redis"`
	S3            map[string]S3Manager     `json:"s3"`
	SFTP          map[string]SFTP          `json:"sftp"`
	Dynamo        map[string]Dynamo        `json:"dynamo"`
}

// Messaging configuration datastructure.
type Messaging struct {
	NSQ     NSQ         `json:"nsq"`
	Kinesis Kinesis     `json:"kinesis"`
	Kafka   Kafka       `json:"kafka"`
	SQS     SQSResource `json:"sqs"`
}

func parse(arn string, path ...string) (string, string, error) {
	if arn == "" {
		return "", "", fmt.Errorf("arn is empty")
	}
	var arp = strings.Split(strings.TrimPrefix(arn, resourcePrefix), "/")
	var subpath = len(arp) - len(path)

	if subpath <= 0 {
		return "", "", fmt.Errorf("mismatch arn and type of resource")
	}

	var i = 0
	var k string
	for i, k = range path {
		if arp[i] != k {
			return "", "", fmt.Errorf("expected arn to contain [%s] at position %d [%+v]", k, i, arp)
		}
	}

	switch subpath {
	case 1:
		return arp[i+1], "", nil
	case 2:
		return arp[i+1], arp[i+2], nil
	default:
		return "", "", fmt.Errorf("arn has wrong path length [%d], expected %d to %d path entries", len(arp), len(path)+1, len(path)+2)
	}
}

// Resource is the base resource which only provides an accessor for detecting resource location errors.
type Resource struct {
	// Tags are used to classify this resource.
	// These can be things like "datacenter:us-west-1", "account:1234567890", "state:discontinued", etc.
	Tags []string `json:"tags"`
	// Params is how you add custom parameters to be used by the connection implementation.
	Params Params `json:"params"`
	err    error
}

func (r Resource) Error() error {
	return r.err
}

// Validate returns true if the resource is valid.
func (r Resource) Validate() error {
	return r.err
}

type Params map[string]interface{}

func (p Params) String(param string) string {
	if len(p) == 0 {
		return ""
	}
	if value, ok := p[param]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func (p Params) Int(param string) int {
	if len(p) == 0 {
		return 0
	}
	if value, ok := p[param]; ok {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return 0
}

func (p Params) Bool(param string) bool {
	if len(p) == 0 {
		return false
	}
	if value, ok := p[param]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

func (p Params) Float64(param string) float64 {
	if len(p) == 0 {
		return 0
	}
	if value, ok := p[param]; ok {
		if f, ok := value.(float64); ok {
			return f
		}
	}
	return 0
}
