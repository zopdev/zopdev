package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/models"
)

type service interface {
	CreateIntegration(ctx *gofr.Context, provider string) (models.Integration, string, error)
	AssumeRoleAndCreateAdmin(ctx *gofr.Context, req *models.AssumeRole) (map[string]string, error)
}
