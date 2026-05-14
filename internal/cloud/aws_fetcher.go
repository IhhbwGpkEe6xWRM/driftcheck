package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceAttributes holds the live attribute map for a cloud resource.
type ResourceAttributes map[string]interface{}

// Fetcher retrieves live resource attributes from AWS.
type Fetcher struct {
	ec2Client *ec2.Client
	s3Client  *s3.Client
}

// NewFetcher creates a Fetcher using the default AWS config chain.
func NewFetcher(ctx context.Context) (*Fetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}
	return &Fetcher{
		ec2Client: ec2.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}, nil
}

// FetchResource returns live attributes for the given resource type and ID.
// Supported types: aws_instance, aws_s3_bucket.
func (f *Fetcher) FetchResource(ctx context.Context, resourceType, resourceID string) (ResourceAttributes, error) {
	switch resourceType {
	case "aws_instance":
		return f.fetchEC2Instance(ctx, resourceID)
	case "aws_s3_bucket":
		return f.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (f *Fetcher) fetchEC2Instance(ctx context.Context, instanceID string) (ResourceAttributes, error) {
	out, err := f.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instance %s: %w", instanceID, err)
	}
	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}
	inst := out.Reservations[0].Instances[0]
	attrs := ResourceAttributes{
		"instance_type": string(inst.InstanceType),
		"instance_state": string(inst.State.Name),
		"ami":           aws.ToString(inst.ImageId),
		"private_ip":    aws.ToString(inst.PrivateIpAddress),
		"public_ip":     aws.ToString(inst.PublicIpAddress),
	}
	return attrs, nil
}

func (f *Fetcher) fetchS3Bucket(ctx context.Context, bucketName string) (ResourceAttributes, error) {
	_, err := f.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("head bucket %s: %w", bucketName, err)
	}
	attrs := ResourceAttributes{
		"bucket": bucketName,
	}
	return attrs, nil
}
