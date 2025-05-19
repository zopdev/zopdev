package models

type Integration struct {
	CloudformationURL string `json:"cloudformation_url"`
	IntegrationID     string `json:"integration_id"`
}
type AssumeRole struct {
	IntegrationID string `json:"integration_id"`
	AccountID     string `json:"account_id"`
}
