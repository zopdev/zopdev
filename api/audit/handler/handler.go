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
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	return h.service.RunAll(ctx, cloudAccId)
}

func (h *Handler) RunById(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	ruleId := strings.TrimSpace(ctx.PathParam("ruleId"))
	if ruleId == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"ruleId"}}
	}

	return h.service.RunById(ctx, ruleId, cloudAccId)
}

func (h *Handler) RunByCategory(ctx *gofr.Context) (any, error) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	cloudAccId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	category := strings.TrimSpace(ctx.PathParam("category"))
	if category == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"category"}}
	}

	return h.service.RunByCategory(ctx, category, cloudAccId)
}
