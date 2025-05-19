package handler

import (
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/integration/models"
)

// Handler handles HTTP requests for AWS integration.
type Handler struct {
	service service
}

// New creates a new Handler instance.
func New(svc service) *Handler {
	return &Handler{service: svc}
}

// CreateIntegration handles GET request to get integration form data or account list.
// It validates the provider and returns integration details with CloudFormation URL.
func (h *Handler) CreateIntegration(ctx *gofr.Context) (any, error) {
	provider := strings.ToLower(strings.TrimSpace(ctx.PathParam("provider")))
	if provider == "" {
		return nil, http.ErrorMissingParam{Params: []string{"provider"}}
	}

	// Validate provider.
	if !isValidProvider(provider) {
		return nil, http.ErrorInvalidParam{Params: []string{"provider"}}
	}

	integration, cfnURL, err := h.service.CreateIntegration(ctx, provider)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"data":               integration,
		"cloudformation_url": cfnURL,
	}, nil
}

// AssumeRole handles POST request to create integration.
// It validates the provider and request body, then creates a temporary admin user.
func (h *Handler) AssumeRole(ctx *gofr.Context) (any, error) {
	provider := strings.ToLower(strings.TrimSpace(ctx.PathParam("provider")))
	if provider == "" {
		return nil, http.ErrorMissingParam{Params: []string{"provider"}}
	}

	// Validate provider
	if !isValidProvider(provider) {
		return nil, http.ErrorInvalidParam{Params: []string{"provider"}}
	}

	var req models.AssumeRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	// Set provider in request
	req.Provider = provider

	return h.service.AssumeRoleAndCreateTemporaryAdmin(ctx, &req)
}

// isValidProvider checks if the provider is supported.
// Currently only AWS is supported.
func isValidProvider(provider string) bool {
	supportedProviders := map[string]bool{
		"aws": true,
	}

	return supportedProviders[provider]
}
