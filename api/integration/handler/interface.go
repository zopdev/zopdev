package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/models"
)

type service interface {
	GetIntegrationURL(ctx *gofr.Context, provider string) (models.AWSIntegrationINFO, error)
	AssumeRoleAndCreateAdmin(ctx *gofr.Context, req *models.RoleRequest) (map[string]string, error)
}
