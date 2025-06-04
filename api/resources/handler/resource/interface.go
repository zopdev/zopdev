package resource

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/service/resource"
)

type Service interface {
	GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]models.Resource, error)
	SyncResources(ctx *gofr.Context, id int64) ([]models.Resource, error)
	ChangeState(ctx *gofr.Context, resDetails resource.ResourceDetails) error
}
