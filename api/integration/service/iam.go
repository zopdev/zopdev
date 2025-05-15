package service

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func CreateAdminUserWithGroup(ctx context.Context, accessKey, secretKey, sessionToken, userName, groupName string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)),
	)
	if err != nil {
		return "", "", errors.New("failed to load AWS config")
	}

	iamClient := iam.NewFromConfig(cfg)

	// Create group if not exists
	_, err = iamClient.CreateGroup(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errors.New("failed to create IAM group")
	}

	// Attach AdministratorAccess policy to group
	_, err = iamClient.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})
	if err != nil {
		return "", "", errors.New("failed to attach AdministratorAccess policy to group")
	}

	// Create user if not exists
	_, err = iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(userName),
	})
	if err != nil && !isEntityAlreadyExists(err) {
		return "", "", errors.New("failed to create IAM user")
	}

	// Add user to group
	_, err = iamClient.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})
	if err != nil {
		return "", "", errors.New("failed to add user to group")
	}

	// Create access key for user
	accessKeyOut, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return "", "", errors.New("failed to create access key for user")
	}

	return *accessKeyOut.AccessKey.AccessKeyId, *accessKeyOut.AccessKey.SecretAccessKey, nil
}

func isEntityAlreadyExists(err error) bool {
	return err != nil && (contains(err.Error(), "EntityAlreadyExists") || contains(err.Error(), "already exists"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
