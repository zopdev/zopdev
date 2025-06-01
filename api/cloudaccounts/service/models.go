package service

type gcpCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type ociCredentials struct {
	TenancyOCID string `json:"tenancy_ocid"`
	UserOCID    string `json:"user_ocid"`
	Region      string `json:"region"`
	Fingerprint string `json:"fingerprint"`
	PrivateKey  string `json:"private_key"`
	Compartment string `json:"compartment"`
}

type awsCredentials struct {
	AccessKey    string `json:"aws_access_key_id"`
	AccessSecret string `json:"aws_secret_access_key"`
}

type DeploymentSpaceOptions struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type AWSIntegrationINFO struct {
	CloudformationURL string `json:"cloudformation_url"`
	IntegrationID     string `json:"integration_id"`
}

type RoleRequest struct {
	CloudAccountName string `json:"cloud_account_name"`
	IntegrationID    string `json:"integration_id"`
	AccountID        string `json:"account_id"`
}
