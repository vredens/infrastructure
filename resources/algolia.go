package resources

import (
	"fmt"
)

// Algolia configuration datastructure.
type Algolia struct {
	Resource
	ApplicationID string `json:"application_id"`
	APIKey        string `json:"api_key"`
	IndexPrefix   string `json:"index_prefix"`
}

func (r Algolia) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.ApplicationID == "" {
		return fmt.Errorf("empty application id")
	}
	if r.APIKey == "" {
		return fmt.Errorf("empty api key")
	}
	if r.IndexPrefix == "" {
		return fmt.Errorf("index prefix is empty")
	}
	return r.err
}

// IndexNameFor computes the final index name for the provided index using the configured index prefix.
func (cfg Algolia) IndexNameFor(index string) string {
	return cfg.IndexPrefix + index
}

func (r Algolia) sanitize() Algolia {
	r.err = r.Validate()
	return r
}

// LocateAlgoliaResource returns an AlgoliaResource definition.
func (irl Locator) LocateAlgoliaResource(arn string) Algolia {
	var resource Algolia
	var found bool
	var name, _, err = parse(arn, "storage", "algolia")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.Algolia[name]; !found {
		resource.Resource.err = fmt.Errorf("algolia resource not found")
		return resource
	}
	return resource.sanitize()
}
