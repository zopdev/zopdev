package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/service"
)

type Service interface {
	GetResources(ctx *gofr.Context, id int64, resources []string) ([]models.Instance, error)
	ChangeState(ctx *gofr.Context, resDetails service.ResourceDetails) error
}
