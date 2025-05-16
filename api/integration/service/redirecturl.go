package service

import (
	"fmt"
	"net/url"

	"github.com/zopdev/zopdev/api/integration/models"
)

func GenerateCloudFormationURL(integration models.Integration, permissionLevel, trustedPrincipalArn string) string {
	region := "us-east-1"

	// Base CloudFormation URL
	baseURL := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudformation/home", region)

	stackName := fmt.Sprintf("Zopdev-%s", integration.IntegrationID)

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
		url.QueryEscape(integration.TemplateURL),
		url.QueryEscape(stackName),
		url.QueryEscape(integration.IntegrationID),
		url.QueryEscape(integration.ExternalID),
		url.QueryEscape(trustedPrincipalArn),
		url.QueryEscape(permissionLevel),
	)

	return cfnURL
}
