package service

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/models"
	"github.com/zopdev/zopdev/api/integration/store"
)

const s3TemplateBaseURL = "https://zopdev-aws-one-click.s3.us-east-1.amazonaws.com/aws-service.yaml"

var (
	errMissingIntegrationOrAccountID = errors.New("missing required fields: integration_id or account_id")
	errFailedToCreateAdminUserGroup  = errors.New("failed to create admin user/group")
)

type IntegrationService struct {
	store               *store.Store
	trustedPrincipalArn string
}

func New(store *store.Store, accountID string) *IntegrationService {
	return &IntegrationService{
		store:               store,
		trustedPrincipalArn: fmt.Sprintf("arn:aws:iam::%s:role/CrossAccountAssumer", accountID),
	}
}

func (s *IntegrationService) CreateIntegration(ctx *gofr.Context, permissionLevel string) (models.Integration, error) {
	integrationID := uuid.New().String()
	externalID := fmt.Sprintf("ext-%s", integrationID)

	integration := models.Integration{
		IntegrationID: integrationID,
		ExternalID:    externalID,
		TemplateURL:   s3TemplateBaseURL,
		RoleName:      fmt.Sprintf("CrossAccountAccessRole-%s", integrationID),
	}

	err := s.store.SaveIntegration(ctx, integration)
	return integration, err
}

func (s *IntegrationService) GetIntegration(ctx *gofr.Context, id string) (models.Integration, error) {
	return s.store.GetIntegration(ctx, id)
}

func (s *IntegrationService) CreateIntegrationWithURL(ctx *gofr.Context, permissionLevel string) (models.Integration, string, error) {
	integration, err := s.CreateIntegration(ctx, permissionLevel)
	if err != nil {
		return integration, "", err
	}
	cfnURL := GenerateCloudFormationURL(integration, permissionLevel, s.trustedPrincipalArn)
	return integration, cfnURL, nil
}

func (s *IntegrationService) AssumeRoleWithOptionalAdminUser(ctx *gofr.Context, integrationID, accountID, userName, groupName string) (map[string]string, error) {
	if integrationID == "" || accountID == "" {
		return nil, errMissingIntegrationOrAccountID
	}

	// Generate user_name and group_name
	randomSuffix := func(n int) string {
		b := make([]byte, n)
		rand.Read(b)
		return fmt.Sprintf("%x", b)[:n]
	}
	if userName == "" {
		suffix := randomSuffix(6)
		userName = "Zop-Admin-" + suffix
	}
	if groupName == "" {
		suffix := randomSuffix(6)
		groupName = "ZopAdminGroup-" + suffix
	}

	integration, err := s.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, integration.RoleName)
	result, err := AssumeRole(roleARN, integration.ExternalID, "session-"+integration.IntegrationID)
	if err != nil {
		return nil, err
	}
	creds := result.Credentials

	ak, sk, err := CreateAdminUserWithGroup(ctx, *creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken, userName, groupName)
	if err != nil {
		return nil, errFailedToCreateAdminUserGroup
	}

	return map[string]string{
		"iam_user_access_key_id":     ak,
		"iam_user_secret_access_key": sk,
	}, nil
}
