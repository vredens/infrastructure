package resources

// Dynamo resource data structure.
type Dynamo struct {
	Resource
	Session AWSSession `json:"session"`
}

// Validate returns true if the resource is valid.
func (r Dynamo) Validate() error {
	if r.err != nil {
		return r.err
	}
	if err := r.Session.Validate(); err != nil {
		return err
	}
	return r.err
}

func (r Dynamo) sanitize() Dynamo {
	r.err = r.Validate()
	return r
}

// LocateDynamoResource ...
func (irl *Locator) LocateDynamoResource(arn string) Dynamo {
	var resource Dynamo
	var found bool
	var name, _, err = parse(arn, "storage", "dynamo")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.Dynamo[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
