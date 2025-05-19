package models

// Integration represents a cloud provider integration.
type Integration struct {
	IntegrationID string `json:"integration_id"`
	ExternalID    string `json:"external_id"`
	TemplateURL   string `json:"template_url"`
	RoleName      string `json:"role_name"`
	Provider      string `json:"provider"`
}

// AssumeRoleRequest represents the request parameters for assuming a role.
type AssumeRoleRequest struct {
	IntegrationID string `json:"integration_id"`
	AccountID     string `json:"account_id"`
	UserName      string `json:"user_name"`
	GroupName     string `json:"group_name"`
	Provider      string `json:"provider"`
}
