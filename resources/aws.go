package resources

// AWSSession defines the configuration for an aws session.
type AWSSession struct {
	Resource
	Endpoint                  string         `json:"endpoint"`
	Region                    string         `json:"region"`
	Role                      string         `json:"role"`
	Credentials               AWSCredentials `json:"credentials"`
	DisableSSL                bool           `json:"disable_ssl"`
	S3ForcePathStyle          bool           `json:"force_path_style"`
	DisableEndpointHostPrefix bool           `json:"disable_endpoint_host_prefix"`
}

func (a AWSSession) sanitize() AWSSession {
	return a
}

func (irl Locator) LocateAWSSession(arn string) AWSSession {
	var resource AWSSession
	var found bool
	var name, _, err = parse(arn, "cloud", "aws")
	if err != nil {
		resource.Resource.err = err
		return resource
	}
	if resource, found = irl.Cloud.AWS[name]; !found {
		resource.Resource.err = ErrResourceNotFound
		return resource
	}
	return resource.sanitize()
}

// AWSCredentials defines the credentials configuration.
type AWSCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Token           string `json:"token"`
}

// IsZero checks if the struct is the zero value.
func (a AWSCredentials) IsZero() bool {
	return a.AccessKeyID == "" && a.SecretAccessKey == "" && a.Token == ""
}
