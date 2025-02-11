package aws

import (
	"encoding/json"
	"fmt"
	"os"
)

// AWS gateway allows retrieving info from the cloud provider if available.
type AWS struct{}

var core AWS

func New() AWS {
	return core
}

// AvailabilityZone tries to get the availability zone your code is running in.
// It returns empty string if none can be found.
//
// Availability zone is available if running from
//   - ECS container instances
//   - (TODO) EC2 instances
func (aws AWS) AvailabilityZone() string {
	meta, err := aws.ECSContainerMetadata()
	if err != nil {
		return ""
	}
	return meta.AvailabilityZone
}

// ECSContainerMetadata information from the file in env var ECS_CONTAINER_METADATA_FILE.
// If there is no value for the env var then an empty metadata is returned.
//
// Check https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-metadata.html.
func (aws AWS) ECSContainerMetadata() (ECSContainerMetadata, error) {
	var filepath = os.Getenv("ECS_CONTAINER_METADATA_FILE")
	var metadata ECSContainerMetadata

	if filepath == "" {
		return metadata, nil
	}

	var payload, err = os.ReadFile(filepath)
	if err != nil {
		return metadata, fmt.Errorf("failed to read ECS container metadata file [%s]; %w", filepath, err)
	}
	err = json.Unmarshal(payload, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("failed to unmarshal ECS metadata file; %w", err)
	}

	return metadata, nil
}
