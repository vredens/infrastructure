package resources

import "fmt"

// SFTP data structure.
type SFTP struct {
	Resource
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Pass       string `json:"pass"`
	PrivateKey struct {
		Value      string `json:"value"`
		Path       string `json:"path"`
		Passphrase string `json:"passphrase"`
	} `json:"private_key"`
	HostKey string `json:"host_key"`
}

// Validate returns true if the resource is valid.
func (r SFTP) Validate() error {
	if r.err != nil {
		return r.err
	}

	if r.Host == "" {
		return fmt.Errorf("no host defined")
	}
	if r.User == "" {
		return fmt.Errorf("no user defined")
	}
	if r.Pass == "" && r.PrivateKey.Path == "" && r.PrivateKey.Value == "" {
		return fmt.Errorf("no pass/private_key defined")
	}
	if r.PrivateKey.Passphrase != "" && r.PrivateKey.Path == "" && r.PrivateKey.Value == "" {
		return fmt.Errorf("passphrase provided but no private key or path to a private key file was defined")
	}
	if r.PrivateKey.Value != "" && r.PrivateKey.Path != "" {
		return fmt.Errorf("private key must be either the value or path to the file containing the private key, can not have both configured")
	}

	return r.err
}

func (r SFTP) sanitize() SFTP {
	if r.Port == 0 {
		r.Port = 22
	}
	r.err = r.Validate()
	return r
}

// LocateSFTPResource pointed to by the ARN.
func (irl Locator) LocateSFTPResource(arn string) SFTP {
	var resource SFTP
	var found bool
	var name, _, err = parse(arn, "storage", "sftp")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Databases.SFTP[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}
