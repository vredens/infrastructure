package resources

import (
	"fmt"
	"net/url"
)

// Postgres configuration datastructure.
type Postgres struct {
	Resource
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	// DSNParams are extra connection parameters to be appended to the DSN in the format of key=value.
	DSNParams map[string]string `json:"dsn_params"`
}

// Validate returns an error if the resource is invalid.
func (r *Postgres) Validate() error {
	if r.err != nil {
		return r.err
	}
	if r.Host == "" {
		return fmt.Errorf("postgres host configuration undefined")
	}
	if r.Database == "" {
		return fmt.Errorf("postgres database configuration undefined")
	}
	if r.User == "" {
		return fmt.Errorf("postgres user configuration undefined")
	}
	return r.err
}

// GetDSN for this configuration. Returns an empty string if configuration is incomplete.
func (cfg Postgres) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.Database)
}

// GetFullDSN which includes extra connection parameters in the format <dsn>?<key>=<value>. Returns an empty string if configuration is incomplete.
func (cfg Postgres) GetFullDSN() string {
	if len(cfg.Params) == 0 {
		return cfg.GetDSN()
	}

	var extra string
	for k, v := range cfg.DSNParams {
		if extra != "" {
			extra += "&"
		}
		extra += fmt.Sprintf("%s=%s", k, v)
	}

	return fmt.Sprintf("%s?%s", cfg.GetDSN(), extra)
}

func (r Postgres) sanitize() Postgres {
	if r.Port == 0 {
		r.Port = 5432
	}
	r.err = r.Validate()
	return r
}

// LocatePostgresResource ...
func (irl Locator) LocatePostgresResource(arn string) Postgres {
	var resource Postgres
	var found bool
	var name, _, err = parse(arn, "storage", "postgres")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.Postgres[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
