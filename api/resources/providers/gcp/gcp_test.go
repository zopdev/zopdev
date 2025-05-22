package gcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
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
