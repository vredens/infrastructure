package resources

import "fmt"

// S3Manager data structure.
type S3Manager struct {
	Resource
	Bucket  string     `json:"bucket"`
	Session AWSSession `json:"session"`
}

// Validate returns true if the resource is valid.
func (r S3Manager) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Bucket == "" {
		return fmt.Errorf("s3 bucket can not be empty")
	}
	if err := r.Session.Validate(); err != nil {
		return err
	}
	return r.err
}

func (r S3Manager) sanitize() S3Manager {
	r.err = r.Validate()
	return r
}

// LocateS3ManagerResource ...
func (irl Locator) LocateS3ManagerResource(arn string) S3Manager {
	var resource S3Manager
	var found bool
	var name, _, err = parse(arn, "storage", "s3")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.S3[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
