package models

type AWSIntegrationINFO struct {
	CloudformationURL string `json:"cloudformation_url"`
	IntegrationID     string `json:"integration_id"`
}

type RoleRequest struct {
	IntegrationID string `json:"integration_id"`
	AccountID     string `json:"account_id"`
}
