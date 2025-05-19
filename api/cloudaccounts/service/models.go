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

type awsCredentials struct {
	AccessKey    string `json:"aws_access_key_id"`
	AccessSecret string `json:"aws_secret_access_key"`
}

type DeploymentSpaceOptions struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}
