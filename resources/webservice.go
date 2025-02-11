package resources

import "fmt"

// Webservice resource configuration datastructure.
type Webservice struct {
	Resource
	BaseURL       string            `json:"url"`
	Headers       map[string]string `json:"headers"`
	Authorisation struct {
		Type string `json:"type"`
		Key  string `json:"key"`
	} `json:"authorisation"`
}

// Validate resource.
func (r Webservice) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.BaseURL == "" {
		return fmt.Errorf("empty url")
	}
	return r.err
}

func (r Webservice) sanitize() Webservice {
	r.err = r.Validate()
	return r
}

// LocateWebserviceResource pointed to by the ARN.
func (irl Locator) LocateWebserviceResource(arn string) Webservice {
	var resource Webservice
	var found bool
	var name, _, err = parse(arn, "webservices")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Webservices[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
