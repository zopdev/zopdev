package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	errFailedToLoadAWSConfig            = errors.New("failed to load AWS config")
	errFailedToCreateIAMGroup           = errors.New("failed to create IAM group")
	errFailedToAttachAdminPolicyToGroup = errors.New("failed to attach AdministratorAccess policy to group")
	errFailedToCreateIAMUser            = errors.New("failed to create IAM user")
	errFailedToAddUserToGroup           = errors.New("failed to add user to group")
	errFailedToCreateAccessKeyForUser   = errors.New("failed to create access key for user")
)

// createAdminUserWithGroup creates an IAM user and group with admin access.
func createAdminUserWithGroup(ctx context.Context,
	accessKey, secretKey, sessionToken, userName, groupName string,
) (accessKeyID, secretAccessKey string, err error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)),
	)
	if err != nil {
		return "", "", errFailedToLoadAWSConfig
	}

	iamClient := iam.NewFromConfig(cfg)

	// Create group if not exists.
	_, err = iamClient.CreateGroup(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errFailedToCreateIAMGroup
	}

	// Attach AdministratorAccess policy to group.
	_, err = iamClient.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String("arn:providerAWS:iam::providerAWS:policy/AdministratorAccess"),
	})
	if err != nil {
		return "", "", errFailedToAttachAdminPolicyToGroup
	}

	// Create user if not exists.
	_, err = iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(userName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errFailedToCreateIAMUser
	}

	// Add user to group.
	_, err = iamClient.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})
	if err != nil {
		return "", "", errFailedToAddUserToGroup
	}

	// Create access key for user.
	accessKeyOut, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return "", "", errFailedToCreateAccessKeyForUser
	}

	return *accessKeyOut.AccessKey.AccessKeyId, *accessKeyOut.AccessKey.SecretAccessKey, nil
}

// isEntityAlreadyExists checks if the error indicates an entity already exists.
func isEntityAlreadyExists(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "EntityAlreadyExists") || strings.Contains(err.Error(), "already exists"))
}

// generateCloudFormationURL generates a CloudFormation console URL for stack creation.
func generateCloudFormationURL(integrationID, externalID, _, permissionLevel, trustedPrincipalArn string) string {
	region := "us-east-1"
	templateURL := s3TemplateBaseURL

	// Base CloudFormation URL.
	baseURL := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudformation/home", region)

	stackName := fmt.Sprintf("Zopdev-%s", integrationID)

	// Quick create stack parameters
	cfnURL := fmt.Sprintf("%s?region=%s#/stacks/quickcreate"+
		"?templateURL=%s"+
		"&stackName=%s"+
		"&param_IntegrationId=%s"+
		"&param_ExternalId=%s"+
		"&param_TrustedPrincipalArn=%s"+
		"&param_PermissionLevel=%s",
		baseURL,
		region,
		url.QueryEscape(templateURL),
		url.QueryEscape(stackName),
		url.QueryEscape(integrationID),
		url.QueryEscape(externalID),
		url.QueryEscape(trustedPrincipalArn),
		url.QueryEscape(permissionLevel),
	)

	return cfnURL
}

// assumeRole assumes an IAM role using AWS STS.
func assumeRole(roleArn, externalID, sessionName string) (*sts.AssumeRoleOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(cfg)

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		ExternalId:      aws.String(externalID),
	}

	return client.AssumeRole(context.TODO(), input)
}
