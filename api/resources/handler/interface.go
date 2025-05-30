package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/service"
)

type Service interface {
	GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]models.Instance, error)
	SyncResources(ctx *gofr.Context, id int64) ([]models.Instance, error)
	ChangeState(ctx *gofr.Context, resDetails service.ResourceDetails) error
}
