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

// testServiceAccountCreds returns a minimal set of service account credentials for testing.
func testServiceAccountCreds() map[string]string {
	return map[string]string{
		"type":           "service_account",
		"project_id":     "test-project",
		"private_key_id": "test-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\n" +
			"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9QFi67K6u2X4K\n" +
			"-----END PRIVATE KEY-----\n",
		"client_email":                "test@test-project.iam.gserviceaccount.com",
		"client_id":                   "test-client-id",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-project.iam.gserviceaccount.com",
	}
}

// setupMockServer creates a mock HTTP server for testing.
func setupMockServer(t *testing.T, response string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(response)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
}

func Test_ListResources_InvalidCredentials(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	got, err := c.ListResources(ctx, map[string]string{}, models.ResourceFilter{
		ResourceTypes: []string{"SQL"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cloud credentials")
	assert.Nil(t, got)
}

func Test_ListResources_UnknownResourceType(t *testing.T) {
	srv := setupMockServer(t, `{"items": []}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	got, err := c.ListResources(ctx, testServiceAccountCreds(), models.ResourceFilter{
		ResourceTypes: []string{"UNKNOWN"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	assert.Nil(t, got)
}

func Test_ListResources_EmptyResourceTypes(t *testing.T) {
	srv := setupMockServer(t, `{"items": []}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	got, err := c.ListResources(ctx, testServiceAccountCreds(), models.ResourceFilter{
		ResourceTypes: []string{},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "private key should be a PEM")
	assert.Nil(t, got)
}

func Test_ListResources_AllResourceType(t *testing.T) {
	srv := setupMockServer(t, `{"items": []}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	got, err := c.ListResources(ctx, testServiceAccountCreds(), models.ResourceFilter{
		ResourceTypes: []string{"ALL"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "private key should be a PEM")
	assert.Nil(t, got)
}

func Test_StartResource_InvalidCredentials(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StartResource(ctx, map[string]string{}, &models.Resource{
		Type: "SQL",
		UID:  "project-id/instance-name",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cloud credentials")
}

func Test_StartResource_UnknownResourceType(t *testing.T) {
	srv := setupMockServer(t, `{"status": "DONE"}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StartResource(ctx, testServiceAccountCreds(), &models.Resource{
		Type: "UNKNOWN",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func Test_StartResource_InvalidUID(t *testing.T) {
	srv := setupMockServer(t, `{"status": "DONE"}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StartResource(ctx, testServiceAccountCreds(), &models.Resource{
		Type: "SQL",
		UID:  "invalid-uid",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameter")
}

func Test_StopResource_InvalidCredentials(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StopResource(ctx, map[string]string{}, &models.Resource{
		Type: "SQL",
		UID:  "project-id/instance-name",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cloud credentials")
}

func Test_StopResource_UnknownResourceType(t *testing.T) {
	srv := setupMockServer(t, `{"status": "DONE"}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StopResource(ctx, testServiceAccountCreds(), &models.Resource{
		Type: "UNKNOWN",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func Test_StopResource_InvalidUID(t *testing.T) {
	srv := setupMockServer(t, `{"status": "DONE"}`)
	defer srv.Close()

	c := &Client{}
	ctx := &gofr.Context{}

	err := c.StopResource(ctx, testServiceAccountCreds(), &models.Resource{
		Type: "SQL",
		UID:  "invalid-uid",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameter")
}
