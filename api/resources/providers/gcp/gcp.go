package gcp

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2/google"

	gmonitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	sql "github.com/zopdev/zopdev/api/resources/gcp/database"
	metric "github.com/zopdev/zopdev/api/resources/gcp/monitoring"
)

type Client struct{}

func New() *Client { return &Client{} }

func (*Client) NewGoogleCredentials(ctx context.Context, cred any, scopes ...string) (*google.Credentials, error) {
	var gcpCreds credentials

	b, _ := json.Marshal(cred)
	if err := json.Unmarshal(b, &gcpCreds); err != nil {
		return nil, err
	}

	creds, err := google.CredentialsFromJSON(ctx, b, scopes...)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func (*Client) NewSQLInstanceLister(ctx context.Context, opts ...option.ClientOption) (InstanceLister, error) {
	admin, err := sqladmin.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &sql.Client{SQL: admin.Instances}, nil
}

func (*Client) NewMetricsClient(ctx context.Context, opts ...option.ClientOption) (TimeSeriesLister, error) {
	mCl, err := gmonitoring.NewMetricClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &metric.Client{MetricClient: mCl}, nil
}
