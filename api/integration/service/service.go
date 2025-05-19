package service

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/models"
)

const (
	s3TemplateBaseURL      = "https://zopdev-aws-one-click.s3.us-east-1.amazonaws.com/aws-service.yaml"
	defaultPermissionLevel = "Admin"
	awsProvider            = "aws"
)

var (
	errMissingIntegrationOrAccountID = errors.New("missing required fields: integration_id or account_id")
	errFailedToCreateAdminUserGroup  = errors.New("failed to create admin user/group")
	errUnsupportedProvider           = errors.New("unsupported provider")
)

type IntegrationService struct {
	trustedPrincipalArn string
}

func New(accountID string) *IntegrationService {
	return &IntegrationService{
		trustedPrincipalArn: fmt.Sprintf("arn:aws:iam::%s:role/CrossAccountAssumer", accountID),
	}
}

func (s *IntegrationService) CreateIntegration(_ *gofr.Context, provider string) (models.Integration, string, error) {
	// Validate provider
	if provider != awsProvider {
		return models.Integration{}, "", errUnsupportedProvider
	}

	integrationID := uuid.New().String()
	externalID := fmt.Sprintf("ext-%s", integrationID)
	roleName := fmt.Sprintf("CrossAccountAccessRole-%s", integrationID)

	// Generate CloudFormation URL
	cfnURL := generateCloudFormationURL(integrationID, externalID, roleName, defaultPermissionLevel, s.trustedPrincipalArn)

	// Construct the integration object for the response
	integration := models.Integration{
		CloudformationURL: cfnURL,
		IntegrationID:     integrationID,
	}

	return integration, cfnURL, nil
}

func (*IntegrationService) AssumeRoleAndCreateAdmin(ctx *gofr.Context, req *models.AssumeRole) (map[string]string, error) {
	if req.IntegrationID == "" || req.AccountID == "" {
		return nil, errMissingIntegrationOrAccountID
	}

	// Generate user_name and group_name
	randomSuffix := func(n int) string {
		b := make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			panic(fmt.Sprintf("failed to read random bytes: %v", err))
		}

		return fmt.Sprintf("%x", b)[:n]
	}

	const suffixLength = 6
	suffix := randomSuffix(suffixLength)
	userName := "Zop-Admin-" + suffix
	groupName := "ZopAdminGroup-" + suffix

	externalID := fmt.Sprintf("ext-%s", req.IntegrationID)
	roleName := fmt.Sprintf("CrossAccountAccessRole-%s", req.IntegrationID)
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", req.AccountID, roleName)

	result, err := AssumeRole(roleARN, externalID, "session-"+req.IntegrationID)
	if err != nil {
		return nil, err
	}

	creds := result.Credentials

	ak, sk, err := createAdminUserWithGroup(ctx, *creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken, userName, groupName)
	if err != nil {
		return nil, errFailedToCreateAdminUserGroup
	}

	return map[string]string{
		"iam_user_access_key_id":     ak,
		"iam_user_secret_access_key": sk,
	}, nil
}
