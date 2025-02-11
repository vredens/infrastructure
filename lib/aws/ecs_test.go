package aws_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vredens/infrastructure/lib/aws"
)

func TestECSMetadata(t *testing.T) {
	var err error
	var md aws.ECSContainerMetadata

	conn := aws.New()

	md, err = conn.ECSContainerMetadata()
	assert.Equal(t, aws.ECSContainerMetadata{}, md)
	assert.Nil(t, err)

	os.Setenv("ECS_CONTAINER_METADATA_FILE", "./testdata/ECS_CONTAINER_METADATA.json")
	md, err = conn.ECSContainerMetadata()
	assert.Equal(t, aws.ECSContainerMetadata{
		Cluster:                "default",
		ContainerInstanceARN:   "arn:aws:ecs:us-west-2:012345678910:container-instance/default/1f73d099-b914-411c-a9ff-81633b7741dd",
		TaskARN:                "arn:aws:ecs:us-west-2:012345678910:task/default/2b88376d-aba3-4950-9ddf-bcb0f388a40c",
		TaskDefinitionFamily:   "console-sample-app-static",
		TaskDefinitionRevision: "1",
		ContainerID:            "aec2557997f4eed9b280c2efd7afccdcedfda4ac399f7480cae870cfc7e163fd",
		ContainerName:          "simple-app",
		DockerContainerName:    "/ecs-console-sample-app-static-1-simple-app-e4e8e495e8baa5de1a00",
		ImageID:                "sha256:2ae34abc2ed0a22e280d17e13f9c01aaf725688b09b7a1525d1a2750e2c0d1de",
		ImageName:              "httpd:2.4",
		PortMappings: []aws.ECSContainerPortMapping{
			{
				ContainerPort: 80,
				HostPort:      80,
				BindIp:        "0.0.0.0",
				Protocol:      "tcp",
			},
		},
		Networks: []aws.ECSContainerNetworks{
			{
				NetworkMode:   "bridge",
				IPv4Addresses: []string{"192.0.2.0"},
			},
		},
		MetadataFileStatus:     "READY",
		AvailabilityZone:       "us-east-1b",
		HostPrivateIPv4Address: "192.0.2.0",
		HostPublicIPv4Address:  "203.0.113.0",
	}, md)
	assert.Nil(t, err)
}
