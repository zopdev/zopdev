package handler

import (
	"strconv"
	"strings"

	"github.com/zopdev/zopdev/api/applications/service"
	"github.com/zopdev/zopdev/api/applications/store"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type Handler struct {
	service service.ApplicationService
}

func New(svc service.ApplicationService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) AddApplication(ctx *gofr.Context) (interface{}, error) {
	application := store.Application{}

	err := ctx.Bind(&application)
	if err != nil {
		ctx.Error(err)
		return nil, http.ErrorInvalidParam{Params: []string{"body"}}
	}

	err = validateApplication(&application)
	if err != nil {
		return nil, err
	}

	res, err := h.service.AddApplication(ctx, &application)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) ListApplications(ctx *gofr.Context) (interface{}, error) {
	applications, err := h.service.FetchAllApplications(ctx)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (h *Handler) GetApplication(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	applicationID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	res, err := h.service.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateApplication(application *store.Application) error {
	application.Name = strings.TrimSpace(application.Name)

	params := []string{}
	if application.Name == "" {
		params = append(params, "name")
	}

	if len(params) > 0 {
		return http.ErrorInvalidParam{Params: params}
	}

	return nil
}
