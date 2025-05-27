package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/service"
	"github.com/zopdev/zopdev/api/resources/store"
)

type Service interface {
	GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]store.Resource, error)
	SyncResources(ctx *gofr.Context, id int64) ([]store.Resource, error)
	ChangeState(ctx *gofr.Context, resDetails service.ResourceDetails) error
}
