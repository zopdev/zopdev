package handler

import "gofr.dev/pkg/gofr"

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Audit(ctx *gofr.Context) (any, error) {
	return nil, nil
}
