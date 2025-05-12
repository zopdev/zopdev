package handler

import (
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
)

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RunAll(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	return h.svc.RunAll(ctx, cloudAccID)
}

func (h *Handler) RunByID(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	ruleID := strings.TrimSpace(ctx.PathParam("ruleId"))

	return h.svc.RunByID(ctx, ruleID, cloudAccID)
}

func (h *Handler) RunByCategory(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	category := strings.ToLower(strings.TrimSpace(ctx.PathParam("category")))

	return h.svc.RunByCategory(ctx, category, cloudAccID)
}

func (h *Handler) GetResultByID(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	ruleID := strings.TrimSpace(ctx.PathParam("ruleId"))

	return h.svc.GetResultByID(ctx, cloudAccID, ruleID)
}

func (h *Handler) GetAllResults(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	return h.svc.GetAllResults(ctx, cloudAccID)
}
