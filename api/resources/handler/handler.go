package handler

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
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
