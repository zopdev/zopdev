package handler

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/service"
)

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetResources(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	accID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	resourceType := ctx.Params("type")

	return h.svc.GetResources(ctx, accID, resourceType)
}

func (h *Handler) GetCloudSQLInstances(ctx *gofr.Context) (any, error) {
	var req service.Request

	err := ctx.Bind(&req)
	if err != nil {
		return nil, err
	}

	return h.svc.GetAllSQLInstances(ctx, req)
}
