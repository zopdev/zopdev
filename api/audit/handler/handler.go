package handler

import (
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
)

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RunAll(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	return h.service.RunAll(ctx, cloudAccID)
}

func (h *Handler) RunByID(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	ruleID := strings.TrimSpace(ctx.PathParam("ruleId"))
	if ruleID == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"ruleId"}}
	}

	return h.service.RunByID(ctx, ruleID, cloudAccID)
}

func (h *Handler) RunByCategory(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	category := strings.ToLower(strings.TrimSpace(ctx.PathParam("category")))
	if category == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"category"}}
	}

	return h.service.RunByCategory(ctx, category, cloudAccID)
}
