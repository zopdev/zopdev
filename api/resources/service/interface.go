package service

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/providers/gcp"
)

type GCPClient interface {
	NewGoogleCredentials(ctx context.Context, cred any, scopes ...string) (*google.Credentials, error)
	NewSQLInstanceLister(ctx context.Context, opts ...option.ClientOption) (gcp.SQLClient, error)
}
