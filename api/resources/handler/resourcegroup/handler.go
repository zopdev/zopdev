package resourcegroup

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
)

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAllResourceGroups(ctx *gofr.Context) (any, error) {
	accID, err := getCloudAccountID(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.svc.GetAllResourceGroups(ctx, accID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) GetResourceGroup(ctx *gofr.Context) (any, error) {
	accID, err := getCloudAccountID(ctx)
	if err != nil {
		return nil, err
	}

	rgID, err := getResourceGroupID(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.svc.GetResourceGroupByID(ctx, accID, rgID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) CreateResourceGroup(ctx *gofr.Context) (any, error) {
	accID, err := getCloudAccountID(ctx)
	if err != nil {
		return nil, err
	}

	var rg models.RGCreate

	err = ctx.Bind(&rg)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	rg.CloudAccountID = accID

	res, err := h.svc.CreateResourceGroup(ctx, &rg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) UpdateResourceGroup(ctx *gofr.Context) (any, error) {
	accID, err := getCloudAccountID(ctx)
	if err != nil {
		return nil, err
	}

	groupID, err := getResourceGroupID(ctx)
	if err != nil {
		return nil, err
	}

	var rg models.RGUpdate

	err = ctx.Bind(&rg)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	rg.CloudAccountID = accID
	rg.ID = groupID

	res, err := h.svc.UpdateResourceGroup(ctx, &rg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) DeleteResourceGroup(ctx *gofr.Context) (any, error) {
	accID, err := getCloudAccountID(ctx)
	if err != nil {
		return nil, err
	}

	rgID, err := getResourceGroupID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.svc.DeleteResourceGroup(ctx, accID, rgID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getResourceGroupID(ctx *gofr.Context) (int64, error) {
	rgIDStr := ctx.PathParam("rgID")
	if rgIDStr == "" {
		return 0, gofrHttp.ErrorMissingParam{Params: []string{"rgId"}}
	}

	rgID, err := strconv.ParseInt(rgIDStr, 10, 64)
	if err != nil {
		return 0, gofrHttp.ErrorInvalidParam{Params: []string{"rgId"}}
	}

	return rgID, nil
}

func getCloudAccountID(ctx *gofr.Context) (int64, error) {
	accIDStr := ctx.PathParam("id")
	if accIDStr == "" {
		return 0, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	accID, err := strconv.ParseInt(accIDStr, 10, 64)
	if err != nil {
		return 0, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	return accID, nil
}
