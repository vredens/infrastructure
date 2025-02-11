package resources

import (
	"fmt"
)

// Redis data structure.
type Redis struct {
	Resource
	SentinelAddresses []string `json:"sentinels" mapstructure:"sentinels"`
	MasterName        string   `json:"master_name" mapstructure:"master_name"`
	Address           string   `json:"address" mapstructure:"address"`
	Password          string   `json:"password" mapstructure:"password"`
	DB                int      `json:"db" mapstructure:"db"`
}

// Validate returns true if the resource is valid.
func (r Redis) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Address == "" && (len(r.SentinelAddresses) == 0 || r.MasterName == "") {
		return fmt.Errorf("no sentinels/master configured or no address specified")
	}
	return r.err
}

func (r Redis) sanitize() Redis {
	r.err = r.Validate()
	return r
}

// LocateRedisResource ...
func (irl Locator) LocateRedisResource(arn string) Redis {
	var resource Redis
	var found bool
	var name, _, err = parse(arn, "storage", "redis")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.Redis[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
