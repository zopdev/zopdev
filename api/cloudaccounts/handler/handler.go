package handler

import (
	"fmt"
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/cloudaccounts/service"
	"github.com/zopdev/zopdev/api/cloudaccounts/store"
)

const (
	nameLength = 255
)

type CloudAccountService interface {
	AddCloudAccount(ctx *gofr.Context, accounts *store.CloudAccount) (*store.CloudAccount, error)
	FetchAllCloudAccounts(ctx *gofr.Context) ([]store.CloudAccount, error)

	FetchDeploymentSpace(ctx *gofr.Context, cloudAccountID int64) (interface{}, error)
	ListNamespaces(ctx *gofr.Context, id int64, clusterName, clusterRegion string) (interface{}, error)
	FetchDeploymentSpaceOptions(ctx *gofr.Context, id int64) ([]service.DeploymentSpaceOptions, error)
	FetchCredentials(ctx *gofr.Context, cloudAccountID int64) (interface{}, error)

	CreateCloudAccountConnection(ctx *gofr.Context, req *service.RoleRequest) (*store.CloudAccount, error)
	GetCloudAccountConnectionInfo(_ *gofr.Context, provider string) (service.AWSIntegrationINFO, error)
	GetStackStatus(ctx *gofr.Context, stackName string) (string, error)
}

type Handler struct {
	service CloudAccountService
}

// New creates a new Handler with the provided CloudAccountService.
func New(clService CloudAccountService) Handler {
	return Handler{service: clService}
}

// AddCloudAccount handles the addition of a new CloudAccount, binding input and validating it.
func (h *Handler) AddCloudAccount(ctx *gofr.Context) (interface{}, error) {
	cloudAccount := store.CloudAccount{}

	err := ctx.Bind(&cloudAccount)
	if err != nil {
		ctx.Error(err.Error())
		return nil, http.ErrorInvalidParam{Params: []string{"body"}}
	}

	err = validateCloudAccount(&cloudAccount)
	if err != nil {
		return nil, err
	}

	resp, err := h.service.AddCloudAccount(ctx, &cloudAccount)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetCredentials(ctx *gofr.Context) (interface{}, error) {
	id := strings.TrimSpace(ctx.PathParam("id"))

	cloudAccountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	credentials, err := h.service.FetchCredentials(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

// ListCloudAccounts retrieves all existing CloudAccounts using the service layer.
func (h *Handler) ListCloudAccounts(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.service.FetchAllCloudAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) ListDeploymentSpace(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	cloudAccountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	res, err := h.service.FetchDeploymentSpace(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) ListNamespaces(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	cloudAccountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	clusterName := strings.TrimSpace(ctx.Param("name"))
	clusterRegion := strings.TrimSpace(ctx.Param("region"))

	if clusterName == "" || clusterRegion == "" {
		return nil, http.ErrorInvalidParam{Params: []string{"cluster"}}
	}

	res, err := h.service.ListNamespaces(ctx, cloudAccountID, clusterName, clusterRegion)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Handler) ListDeploymentSpaceOptions(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	cloudAccountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	res, err := h.service.FetchDeploymentSpaceOptions(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// validateCloudAccount checks the required fields and values in a CloudAccount.
func validateCloudAccount(cloudAccount *store.CloudAccount) error {
	params := []string{}

	cloudAccount.Name = strings.TrimSpace(cloudAccount.Name)

	if cloudAccount.Name == "" {
		params = append(params, "name")
	}

	if cloudAccount.Provider == "" {
		params = append(params, "provider")
	}

	if len(params) != 0 {
		return http.ErrorMissingParam{Params: params}
	}

	if len(cloudAccount.Name) > nameLength {
		return http.ErrorInvalidParam{Params: []string{"name"}}
	}

	return nil
}

// GetCloudAccountConnectionInfo handles GET request to get integration form data or account list.
// It validates the provider and returns integration details with CloudFormation URL.
func (h *Handler) GetCloudAccountConnectionInfo(ctx *gofr.Context) (any, error) {
	provider := strings.ToLower(strings.TrimSpace(ctx.PathParam("provider")))
	if provider == "" {
		return nil, http.ErrorMissingParam{Params: []string{"provider"}}
	}

	// Validate provider.
	if !isValidProviderForAutomaticIntegration(provider) {
		return nil, http.ErrorInvalidParam{Params: []string{"provider"}}
	}

	integration, err := h.service.GetCloudAccountConnectionInfo(ctx, provider)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

// CreateCloudAccountConnection handles POST request to create integration.
// It validates the request body and creates a temporary admin user.
func (h *Handler) CreateCloudAccountConnection(ctx *gofr.Context) (any, error) {
	var req service.RoleRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.CreateCloudAccountConnection(ctx, &req)
}

// isValidProviderForAutomaticIntegration checks if the provider is supported.
// Currently only AWS is supported.
func isValidProviderForAutomaticIntegration(provider string) bool {
	supportedProviders := map[string]bool{
		"aws": true,
	}

	return supportedProviders[provider]
}

// GetStackStatus returns the status of a CloudFormation stack.
func (h *Handler) GetStackStatus(ctx *gofr.Context) (interface{}, error) {
	integrationID := ctx.PathParam("integrationId")
	if integrationID == "" {
		return nil, http.ErrorMissingParam{Params: []string{"integrationId"}}
	}

	stackName := fmt.Sprintf("Zopdev-%s", integrationID)

	status, err := h.service.GetStackStatus(ctx, stackName)
	if err != nil {
		return nil, err
	}

	return map[string]string{"status": status}, nil
}
