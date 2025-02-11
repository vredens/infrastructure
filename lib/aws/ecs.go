package aws

// ECSContainerMetadata as per https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-metadata.html.
type ECSContainerMetadata struct {
	Cluster                string
	ContainerInstanceARN   string
	TaskARN                string
	TaskDefinitionFamily   string
	TaskDefinitionRevision string
	ContainerID            string
	ContainerName          string
	DockerContainerName    string
	ImageID                string
	ImageName              string
	PortMappings           []ECSContainerPortMapping
	Networks               []ECSContainerNetworks
	MetadataFileStatus     string
	AvailabilityZone       string
	HostPrivateIPv4Address string
	HostPublicIPv4Address  string
}

func (meta *ECSContainerMetadata) IsZero() bool {
	return meta == nil || meta.ContainerID == ""
}

type ECSContainerPortMapping struct {
	ContainerPort int
	HostPort      int
	BindIp        string
	Protocol      string
}

type ECSContainerNetworks struct {
	NetworkMode   string
	IPv4Addresses []string
}
