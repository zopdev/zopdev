package gcp

import (
	"time"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type SQLClient interface {
	InstanceLister
}

type MetricsClient interface {
	TimeSeriesLister
}

type InstanceLister interface {
	GetAllInstances(ctx *gofr.Context, projectID string) ([]models.SQLInstance, error)
}

type TimeSeriesLister interface {
	GetTimeSeries(ctx *gofr.Context, start, end time.Time, projectID, filter string) ([]models.Metric, error)
}
