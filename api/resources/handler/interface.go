package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Service interface {
	GetResources(ctx *gofr.Context, id int64, resources []string) ([]models.Instance, error)
}
