package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/integration/service"
)

type Handler struct {
	service *service.IntegrationService
}

func New(svc *service.IntegrationService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) CreateIntegration(ctx *gofr.Context) (any, error) {
	permissionLevel := "Admin"
	integration, cfnURL, err := h.service.CreateIntegrationWithURL(ctx, permissionLevel)

	if err != nil {
		return nil, err
	}

	return map[string]any{
		"data":               integration,
		"cloudformation_url": cfnURL,
	}, nil
}

func (h *Handler) AssumeRole(ctx *gofr.Context) (any, error) {
	var req struct {
		IntegrationID string `json:"integration_id"`
		AccountID     string `json:"account_id"`
		UserName      string `json:"user_name"`
		GroupName     string `json:"group_name"`
	}

	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.AssumeRoleWithOptionalAdminUser(ctx, req.IntegrationID, req.AccountID, req.UserName, req.GroupName)
}
