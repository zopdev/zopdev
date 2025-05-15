package gcp

import (
	"context"

	gmonitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	sql "github.com/zopdev/zopdev/api/resources/gcp/database"
	metric "github.com/zopdev/zopdev/api/resources/gcp/monitoring"
)

type Client struct{}

func New() *Client { return &Client{} }

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
