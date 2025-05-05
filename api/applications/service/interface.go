package service

import (
	"github.com/zopdev/zopdev/api/applications/store"
	"gofr.dev/pkg/gofr"
)

type ApplicationService interface {
	AddApplication(ctx *gofr.Context, application *store.Application) (*store.Application, error)
	FetchAllApplications(ctx *gofr.Context) ([]store.Application, error)
	GetApplication(ctx *gofr.Context, id int) (*store.Application, error)
}
