package aws

import (
	"context"
	"encoding/json"
	"errors"

	"gofr.dev/pkg/gofr"

	"gofr.dev/pkg/gofr/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/providers/aws/database"
	"github.com/zopdev/zopdev/api/resources/providers/aws/vm"
)

var (
	ErrInvalidCredentials = errors.New("invalid cloud credentials")
	ErrInitializingClient = errors.New("error initializing AWS client")
	ErrNotImplemented     = errors.New("not implemented")
)

const allResource = "ALL"
const resourceEC2 = "EC2"
const databaseRDS = "RDS"

type Client struct {
}

func New() *Client {
	return &Client{}
}

// NewRDSClient creates a new RDS client with stored credentials.
func (c *Client) NewRDSClient(_ context.Context, creds any) (*database.Client, error) {
	awsCreds, err := getAWSCredentials(creds)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	sess, err := c.newSession(awsCreds.AccessKey, awsCreds.AccessSecret)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &database.Client{RDS: rds.New(sess)}, nil
}

type awsCredentials struct {
	AccessKey    string `json:"aws_access_key_id"`
	AccessSecret string `json:"aws_secret_access_key"`
}

func getAWSCredentials(creds any) (awsCredentials, error) {
	var awsCred awsCredentials

	awsCredBody, _ := json.Marshal(creds)

	err := json.Unmarshal(awsCredBody, &awsCred)
	if err != nil {
		return awsCred, err
	}

	if awsCred.AccessKey == "" && awsCred.AccessSecret == "" {
		return awsCred, http.ErrorMissingParam{Params: []string{"AWSAccessKeyID", "AWSecretAccessKey"}}
	}

	if awsCred.AccessKey == "" {
		return awsCred, http.ErrorMissingParam{Params: []string{"AWSAccessKeyID"}}
	}

	if awsCred.AccessSecret == "" {
		return awsCred, http.ErrorMissingParam{Params: []string{"AWSecretAccessKey"}}
	}

	return awsCred, nil
}

// NewEC2Client creates a new EC2 client with stored credentials.
func (c *Client) NewEC2Client(_ context.Context, creds any) (*vm.Client, error) {
	awsCreds, err := getAWSCredentials(creds)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	sess, err := c.newSession(awsCreds.AccessKey, awsCreds.AccessSecret)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &vm.Client{EC2: ec2.New(sess)}, nil
}

// newSession creates a new AWS session with the stored config and credentials.
func (*Client) newSession(accessKey, secretKey string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"), // Default region, can be overridden
	})
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return sess, nil
}

// shouldIncludeResourceType determines if a resource type should be included based on the filter.
func shouldIncludeResourceType(resourceTypes []string, targetType string) bool {
	if len(resourceTypes) == 0 {
		return true
	}

	for _, t := range resourceTypes {
		if t == targetType || t == allResource {
			return true
		}
	}

	return false
}

// getEC2Resources retrieves EC2 resources using the provided credentials.
func (c *Client) getEC2Resources(ctx *gofr.Context, creds any) ([]models.Resource, error) {
	ec2Client, err := c.NewEC2Client(ctx, creds)
	if err != nil {
		return nil, err
	}

	return ec2Client.GetAllInstances(ctx)
}

// getRDSResources retrieves RDS resources using the provided credentials.
func (c *Client) getRDSResources(ctx *gofr.Context, creds any) ([]models.Resource, error) {
	rdsClient, err := c.NewRDSClient(ctx, creds)
	if err != nil {
		return nil, err
	}

	return rdsClient.GetAllInstances(ctx)
}

// ListResources implements the CloudResourceProvider interface to list AWS resources.
func (c *Client) ListResources(ctx *gofr.Context, creds any, filter models.ResourceFilter) ([]models.Resource, error) {
	var allResources []models.Resource

	includeEC2 := shouldIncludeResourceType(filter.ResourceTypes, resourceEC2)
	includeRDS := shouldIncludeResourceType(filter.ResourceTypes, databaseRDS)

	if includeEC2 {
		instances, err := c.getEC2Resources(ctx, creds)
		if err != nil {
			return nil, err
		}

		allResources = append(allResources, instances...)
	}

	if includeRDS {
		instances, err := c.getRDSResources(ctx, creds)
		if err != nil {
			return nil, err
		}

		allResources = append(allResources, instances...)
	}

	if len(allResources) == 0 {
		return nil, ErrNotImplemented // or a more descriptive error
	}

	return allResources, nil
}

func (c *Client) StartResource(ctx *gofr.Context, creds any, resource *models.Resource) error {
	switch resource.Type {
	case resourceEC2:
		ec2Client, err := c.NewEC2Client(ctx, creds)
		if err != nil {
			return err
		}

		return ec2Client.StartInstance(ctx, resource.UID)
	case databaseRDS:
		rdsClient, err := c.NewRDSClient(ctx, creds)
		if err != nil {
			return err
		}

		return rdsClient.StartInstance(ctx, resource)
	default:
		return ErrNotImplemented
	}
}

func (c *Client) StopResource(ctx *gofr.Context, creds any, resource *models.Resource) error {
	switch resource.Type {
	case "EC2":
		ec2Client, err := c.NewEC2Client(ctx, creds)
		if err != nil {
			return err
		}

		return ec2Client.StopInstance(ctx, resource.UID)
	case "RDS":
		rdsClient, err := c.NewRDSClient(ctx, creds)
		if err != nil {
			return err
		}

		return rdsClient.StopInstance(ctx, resource)
	default:
		return ErrNotImplemented
	}
}
