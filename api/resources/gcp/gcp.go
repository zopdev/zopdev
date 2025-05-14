package gcp

import (
	"context"

	gmonitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	sql "github.com/zopdev/zopdev/api/resources/gcp/database"
	metric "github.com/zopdev/zopdev/api/resources/gcp/monitoring"
)

type client struct{}

func New() *client { return &client{} }

func (*client) NewSQLInstanceLister(ctx context.Context, opts ...option.ClientOption) (InstanceLister, error) {
	admin, err := sqladmin.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &sql.Client{SQL: admin.Instances}, nil
}

func (*client) NewMetricsClient(ctx context.Context, opts ...option.ClientOption) (MetricsClient, error) {
	mCl, err := gmonitoring.NewMetricClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &metric.Client{MetricClient: mCl}, nil
}
