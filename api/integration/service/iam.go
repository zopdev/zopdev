package service

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

var (
	errFailedToLoadAWSConfig            = errors.New("failed to load AWS config")
	errFailedToCreateIAMGroup           = errors.New("failed to create IAM group")
	errFailedToAttachAdminPolicyToGroup = errors.New("failed to attach AdministratorAccess policy to group")
	errFailedToCreateIAMUser            = errors.New("failed to create IAM user")
	errFailedToAddUserToGroup           = errors.New("failed to add user to group")
	errFailedToCreateAccessKeyForUser   = errors.New("failed to create access key for user")
)

func CreateAdminUserWithGroup(ctx context.Context,
	accessKey, secretKey, sessionToken, userName, groupName string,
) (accessKeyID, secretAccessKey string, err error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)),
	)
	if err != nil {
		return "", "", errFailedToLoadAWSConfig
	}

	iamClient := iam.NewFromConfig(cfg)

	// Create group if not exists
	_, err = iamClient.CreateGroup(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errFailedToCreateIAMGroup
	}

	// Attach AdministratorAccess policy to group
	_, err = iamClient.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})
	if err != nil {
		return "", "", errFailedToAttachAdminPolicyToGroup
	}

	// Create user if not exists
	_, err = iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(userName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errFailedToCreateIAMUser
	}

	// Add user to group
	_, err = iamClient.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})
	if err != nil {
		return "", "", errFailedToAddUserToGroup
	}

	// Create access key for user
	accessKeyOut, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return "", "", errFailedToCreateAccessKeyForUser
	}

	return *accessKeyOut.AccessKey.AccessKeyId, *accessKeyOut.AccessKey.SecretAccessKey, nil
}

func isEntityAlreadyExists(err error) bool {
	return err != nil && (contains(err.Error(), "EntityAlreadyExists") || contains(err.Error(), "already exists"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
