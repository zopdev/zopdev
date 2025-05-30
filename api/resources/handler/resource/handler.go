package resource

import (
	"github.com/zopdev/zopdev/api/resources/service/resource"
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

	res, err := h.svc.GetAll(ctx, accID, resourceType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) ChangeState(ctx *gofr.Context) (any, error) {
	var resDetails resource.ResourceDetails

	err := ctx.Bind(&resDetails)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"request body"}}
	}

	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	accID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	resDetails.CloudAccID = accID

	err = h.svc.ChangeState(ctx, resDetails)
	if err != nil {
		return nil, err
	}

	return resDetails, nil
}

func (h *Handler) SyncResources(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	accID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	res, err := h.svc.SyncResources(ctx, accID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
