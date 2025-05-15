package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/service"
)

type Service interface {
	GetAllSQLInstances(ctx *gofr.Context, req service.Request) ([]models.SQLInstance, error)
}
