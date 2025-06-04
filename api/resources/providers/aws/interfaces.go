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

// TODO start and stop instance should be defined here, which would be same for AWS, GCP etc - GCP not being here.
// For testing purpose we are keeping the interfaces here
