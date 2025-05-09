// Package handler provides HTTP handlers for managing deployment spaces.
// It acts as a layer connecting the HTTP interface to the service layer
// for deployment space operations.
package handler

import (
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/deploymentspace/service"
)

// Handler is responsible for handling HTTP requests related to deployment spaces.
// It utilizes the DeploymentSpaceService to perform business logic operations.
type Handler struct {
	service service.DeploymentSpaceService
}

// New initializes a new Handler with the provided DeploymentSpaceService.
//
// Parameters:
//   - svc: An instance of DeploymentSpaceService to handle deployment space operations.
//
// Returns:
//   - An initialized Handler instance.
func New(svc service.DeploymentSpaceService) Handler {
	return Handler{
		service: svc,
	}
}

// Add handles HTTP POST requests to add a new deployment space.
func (h *Handler) Add(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	environmentID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	deploymentSpace := service.DeploymentSpace{}

	err = ctx.Bind(&deploymentSpace)
	if err != nil {
		return nil, err
	}

	err = validate(&deploymentSpace)
	if err != nil {
		return nil, err
	}

	resp, err := h.service.Add(ctx, &deploymentSpace, environmentID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func validate(deploymentSpace *service.DeploymentSpace) error {
	params := []string{}

	if deploymentSpace.CloudAccount.ID == 0 {
		params = append(params, "cloudAccount ID")
	}

	if deploymentSpace.Type.Name == "" {
		params = append(params, "type")
	}

	if len(params) > 0 {
		return http.ErrorMissingParam{Params: params}
	}

	return nil
}

func (h *Handler) ListServices(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	environmentID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	resp, err := h.service.GetServices(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) ListDeployments(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	environmentID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	resp, err := h.service.GetDeployments(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) ListPods(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	environmentID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	resp, err := h.service.GetPods(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) ListCronJobs(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	environmentID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	resp, err := h.service.GetCronJobs(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetService(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	name := ctx.PathParam("name")

	envID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	if name == "" {
		return nil, http.ErrorMissingParam{Params: []string{"name"}}
	}

	resp, err := h.service.GetServiceByName(ctx, envID, name)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetDeployment(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	name := ctx.PathParam("name")

	envID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	if name == "" {
		return nil, http.ErrorMissingParam{Params: []string{"name"}}
	}

	resp, err := h.service.GetDeploymentByName(ctx, envID, name)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetPod(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	name := ctx.PathParam("name")

	envID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	if name == "" {
		return nil, http.ErrorMissingParam{Params: []string{"name"}}
	}

	resp, err := h.service.GetPodByName(ctx, envID, name)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetCronJob(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	name := ctx.PathParam("name")

	envID, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(err, "failed to convert environment id to int")

		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	if name == "" {
		return nil, http.ErrorMissingParam{Params: []string{"name"}}
	}

	resp, err := h.service.GetCronJobByName(ctx, envID, name)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
