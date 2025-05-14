package gcp

import (
	"time"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

type InstanceLister interface {
	GetAllInstances(ctx *gofr.Context, projectId string) ([]models.SQLInstance, error)
}

type MetricsClient interface {
	GetTimeSeries(ctx *gofr.Context, start, end time.Time, projectId, filter string) ([]models.Metric, error)
}
