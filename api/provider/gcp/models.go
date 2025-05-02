package gcp

// gcpCredentials holds the authentication details for a Google Cloud Platform (GCP) account.
// It contains all the necessary fields to authenticate and interact with GCP resources, including
// project ID, private key, client email, and other credentials for OAuth2-based authentication.
//
// This struct is typically used for passing credentials to services that need to authenticate
// against GCP, such as the GCP provider service.

type gcpCredentials struct {
	// Type represents the type of the credentials, typically "service_account".
	Type string `json:"type"`

	// ProjectID is the ID of the GCP project associated with the credentials.
	ProjectID string `json:"project_id"`

	// PrivateKeyID is the identifier for the private key used in the authentication process.
	PrivateKeyID string `json:"private_key_id"`

	// PrivateKey contains the private key used for authentication with the GCP service.
	PrivateKey string `json:"private_key"`

	// ClientEmail is the email address associated with the GCP service account.
	ClientEmail string `json:"client_email"`

	// ClientID is the identifier for the GCP service account client.
	ClientID string `json:"client_id"`

	// AuthURI is the URI for the authorization server used for OAuth2 authentication.
	AuthURI string `json:"auth_uri"`

	// TokenURI is the URI used to obtain the access token for authentication.
	TokenURI string `json:"token_uri"`

	// AuthProviderX509CertURL is the URL of the X.509 certificate used to verify the identity of the authentication provider.
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`

	// ClientX509CertURL is the URL of the X.509 certificate for the client.
	ClientX509CertURL string `json:"client_x509_cert_url"`

	// UniverseDomain represents the domain for the GCP service account.
	UniverseDomain string `json:"universe_domain"`
}
