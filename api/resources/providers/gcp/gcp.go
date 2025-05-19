package gcp

import (
	"context"
	"encoding/json"
	"errors"

	gmonitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	sql "github.com/zopdev/zopdev/api/resources/providers/gcp/database"
	metric "github.com/zopdev/zopdev/api/resources/providers/gcp/monitoring"
)

var (
	ErrInvalidCredentials = errors.New("invalid cloud credentials")
	ErrInitializingClient = errors.New("error initializing client")
)

type Client struct{}

func New() *Client { return &Client{} }

func (*Client) NewGoogleCredentials(ctx context.Context, cred any, scopes ...string) (*google.Credentials, error) {
	var gcpCreds credentials

	b, _ := json.Marshal(cred)
	if err := json.Unmarshal(b, &gcpCreds); err != nil {
		return nil, ErrInvalidCredentials
	}

	creds, err := google.CredentialsFromJSON(ctx, b, scopes...)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return creds, nil
}

func (*Client) NewSQLInstanceLister(ctx context.Context, opts ...option.ClientOption) (SQLClient, error) {
	admin, err := sqladmin.NewService(ctx, opts...)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &sql.Client{SQL: admin.Instances}, nil
}

func (*Client) NewMetricsClient(ctx context.Context, opts ...option.ClientOption) (MetricsClient, error) {
	mCl, err := gmonitoring.NewMetricClient(ctx, opts...)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &metric.Client{MetricClient: mCl}, nil
}
