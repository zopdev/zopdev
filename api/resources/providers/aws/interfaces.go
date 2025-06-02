package aws

import (
	"context"

	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
)

type ResourceClient interface {
	InstanceLister
	Idler
}

type InstanceLister interface {
	GetAllInstances(ctx context.Context) ([]models.Resource, error)
}

type Idler interface {
	StartInstance(ctx *gofr.Context, instanceName string) error
	StopInstance(ctx *gofr.Context, instanceName string) error
}
