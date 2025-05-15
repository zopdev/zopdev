package handler

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/service"
)

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetCloudSQLInstances(ctx *gofr.Context) (any, error) {
	var req service.Request

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return h.svc.GetAllSQLInstances(ctx, req)
}
