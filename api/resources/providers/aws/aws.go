package aws

import (
	"context"
	"encoding/json"
	"errors"

	"gofr.dev/pkg/gofr/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/zopdev/zopdev/api/resources/providers/aws/database"
	"github.com/zopdev/zopdev/api/resources/providers/aws/vm"
)

var (
	ErrInvalidCredentials = errors.New("invalid cloud credentials")
	ErrInitializingClient = errors.New("error initializing AWS client")
)

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
