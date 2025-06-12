package gcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestClient_NewGoogleCredentials(t *testing.T) {
	var cred any

	ctx := context.Background()
	cred = map[string]string{
		"type":       "service_account",
		"project_id": "test-project",
	}

	c := New()
	creds, err := c.NewGoogleCredentials(ctx, cred, "https://www.googleapis.com/auth/cloud-platform")

	assert.NotNil(t, creds)
	require.NoError(t, err)

	cred = `{"type":"service_account","project_id":"test-project"}`
	creds, err = c.NewGoogleCredentials(ctx, cred)

	assert.Nil(t, creds)
	require.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestClient_NewSQLInstanceLister(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))

	ctx := context.Background()
	c := New()
	sql, err := c.NewSQLClient(ctx, option.WithEndpoint(srv.URL), option.WithoutAuthentication())

	require.NoError(t, err)
	assert.NotNil(t, sql)

	sql, err = c.NewSQLClient(ctx, option.WithoutAuthentication(), option.WithCredentialsFile("test.json"))

	assert.Nil(t, sql)
	require.Error(t, err)
	assert.Equal(t, ErrInitializingClient, err)
}

func TestClient_NewMetricsClient(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	c := New()
	met, err := c.NewMetricsClient(ctx, option.WithEndpoint(srv.URL), option.WithoutAuthentication())

	require.NoError(t, err)
	assert.NotNil(t, met)

	met, err = c.NewMetricsClient(ctx, option.WithoutAuthentication(), option.WithCredentialsFile("test.json"))

	assert.Nil(t, met)
	require.Error(t, err)
	assert.Equal(t, ErrInitializingClient, err)
}

func Test_ListResources(t *testing.T) {
	// Create a mock server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items": []}`))
	}))
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		filter       models.ResourceFilter
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			filter: models.ResourceFilter{
				ResourceTypes: []string{"SQL"},
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name: "valid credentials but no resources",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			filter: models.ResourceFilter{
				ResourceTypes: []string{"UNKNOWN"},
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
		{
			name: "empty resource types should include SQL",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			filter: models.ResourceFilter{
				ResourceTypes: []string{},
			},
			wantErr:      true,
			errorMessage: "private key should be a PEM",
		},
		{
			name: "ALL resource type should include SQL",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			filter: models.ResourceFilter{
				ResourceTypes: []string{"ALL"},
			},
			wantErr:      true,
			errorMessage: "private key should be a PEM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListResources(ctx, tt.creds, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Nil(t, got)
			}
		})
	}
}

func Test_StartResource(t *testing.T) {
	// Create a mock server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "DONE"}`))
	}))
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		resource     *models.Resource
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			resource: &models.Resource{
				Type: "SQL",
				UID:  "project-id/instance-name",
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name: "unsupported resource type",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			resource: &models.Resource{
				Type: "UNKNOWN",
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
		{
			name: "invalid resource UID format",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			resource: &models.Resource{
				Type: "SQL",
				UID:  "invalid-uid",
			},
			wantErr:      true,
			errorMessage: "invalid parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.StartResource(ctx, tt.creds, tt.resource)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			}
		})
	}
}

func Test_StopResource(t *testing.T) {
	// Create a mock server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "DONE"}`))
	}))
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		resource     *models.Resource
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			resource: &models.Resource{
				Type: "SQL",
				UID:  "project-id/instance-name",
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name: "unsupported resource type",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			resource: &models.Resource{
				Type: "UNKNOWN",
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
		{
			name: "invalid resource UID format",
			creds: map[string]string{
				"type":                        "service_account",
				"project_id":                  "test-project",
				"private_key_id":              "test-key-id",
				"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n-----END PRIVATE KEY-----\n",
				"client_email":                "test@test-project.iam.gserviceaccount.com",
				"client_id":                   "test-client-id",
				"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
				"token_uri":                   "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
			},
			resource: &models.Resource{
				Type: "SQL",
				UID:  "invalid-uid",
			},
			wantErr:      true,
			errorMessage: "invalid parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.StopResource(ctx, tt.creds, tt.resource)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			}
		})
	}
}
