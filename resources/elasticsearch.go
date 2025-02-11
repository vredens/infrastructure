package resources

import "fmt"

// Elasticsearch resource configuration datastructure.
type Elasticsearch struct {
	Resource
	Hosts       []string `json:"hosts"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	IndexPrefix string   `json:"index_prefix"`
}

// Validate returns true if the resource is valid.
func (r Elasticsearch) Validate() error {
	if r.err != nil {
		return r.err
	}
	if len(r.Hosts) == 0 {
		return fmt.Errorf("no hosts configured")
	}
	return r.err
}

func (r Elasticsearch) sanitize() Elasticsearch {
	r.err = r.Validate()
	return r
}

// LocateElasticResource returns an ElasticsearchResource definition.
func (irl *Locator) LocateElasticResource(arn string) Elasticsearch {
	var resource Elasticsearch
	var found bool
	var name, _, err = parse(arn, "storage", "elasticsearch")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.Elasticsearch[name]; !found {
		resource.Resource.err = fmt.Errorf("elasticsearch resource not found")
		return resource
	}
	return resource.sanitize()
}
