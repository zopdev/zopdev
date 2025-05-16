package models

type Integration struct {
	IntegrationID string `json:"integration_id"`
	ExternalID    string `json:"external_id"`
	TemplateURL   string `json:"template_url"`
	RoleName      string `json:"role_name"`
}
