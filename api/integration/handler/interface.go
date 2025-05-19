package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/models"
)

type service interface {
	CreateIntegration(ctx *gofr.Context, provider string) (models.Integration, string, error)
	AssumeRoleAndCreateTemporaryAdmin(ctx *gofr.Context, req *models.AssumeRoleRequest) (map[string]string, error)
}
