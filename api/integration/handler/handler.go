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

// GetIntegration handles GET request to get integration form data or account list.
// It validates the provider and returns integration details with CloudFormation URL.
func (h *Handler) GetIntegration(ctx *gofr.Context) (any, error) {
	provider := strings.ToLower(strings.TrimSpace(ctx.PathParam("provider")))
	if provider == "" {
		return nil, http.ErrorMissingParam{Params: []string{"provider"}}
	}

	// Validate provider.
	if !isValidProvider(provider) {
		return nil, http.ErrorInvalidParam{Params: []string{"provider"}}
	}

	integration, err := h.service.GetIntegrationURL(ctx, provider)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

// Connect handles POST request to create integration.
// It validates the request body and creates a temporary admin user.
func (h *Handler) Connect(ctx *gofr.Context) (any, error) {
	var req models.RoleRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.AssumeRoleAndCreateAdmin(ctx, &req)
}

// isValidProvider checks if the provider is supported.
// Currently only AWS is supported.
func isValidProvider(provider string) bool {
	supportedProviders := map[string]bool{
		"aws": true,
	}

	return supportedProviders[provider]
}
