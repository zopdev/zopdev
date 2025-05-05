package handler

import (
	"strconv"
	"strings"

	"github.com/zopdev/zopdev/api/environments/service"
	"github.com/zopdev/zopdev/api/environments/store"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

// Handler is responsible for handling HTTP requests related to environments.
// It interacts with the EnvironmentService to perform operations on environments.
type Handler struct {
	service service.EnvironmentService
}

// New creates a new instance of Handler.
//
// Parameters:
//   - svc: The EnvironmentService implementation for handling business logic.
//
// Returns:
//   - *Handler: A new handler instance.
func New(svc service.EnvironmentService) *Handler {
	return &Handler{
		service: svc,
	}
}

// Add handles the HTTP request to add a new environment.
//
// The method extracts the application ID from the request's path parameter, validates the input,
// and delegates the creation of the environment to the service layer.
//
// Parameters:
//   - ctx: The HTTP context, which includes the request and response data.
//
// Returns:
//   - interface{}: The newly created environment record.
//   - error: An error if the operation fails.
func (h *Handler) Add(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	applicationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	environment := store.Environment{}

	err = ctx.Bind(&environment)
	if err != nil {
		return nil, err
	}

	environment.ApplicationID = applicationID

	err = validateEnvironment(&environment)
	if err != nil {
		return nil, err
	}

	res, err := h.service.Add(ctx, &environment)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// List handles the HTTP request to retrieve all environments for a specific application.
//
// The method extracts the application ID from the request's path parameter and
// delegates the fetching of environments to the service layer.
//
// Parameters:
//   - ctx: The HTTP context, which includes the request and response data.
//
// Returns:
//   - interface{}: A slice of environments associated with the application.
//   - error: An error if the operation fails.
func (h *Handler) List(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	id = strings.TrimSpace(id)

	applicationID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	res, err := h.service.FetchAll(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update handles the HTTP request to update multiple environments.
//
// The method binds the input from the request body, validates the data, and
// delegates the update operation to the service layer.
//
// Parameters:
//   - ctx: The HTTP context, which includes the request and response data.
//
// Returns:
//   - interface{}: A slice of updated environment records.
//   - error: An error if the operation fails.
func (h *Handler) Update(ctx *gofr.Context) (interface{}, error) {
	environments := []store.Environment{}

	err := ctx.Bind(&environments)
	if err != nil {
		return nil, err
	}

	for i := range environments {
		err = validateEnvironment(&environments[i])
		if err != nil {
			return nil, err
		}
	}

	res, err := h.service.Update(ctx, environments)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// validateEnvironment validates the input for an environment.
//
// The method checks if required fields are present and trims extra spaces.
//
// Parameters:
//   - environment: A pointer to the Environment struct to be validated.
//
// Returns:
//   - error: An error if validation fails, specifying the missing fields.
func validateEnvironment(environment *store.Environment) error {
	environment.Name = strings.TrimSpace(environment.Name)
	params := []string{}

	if environment.Name == "" {
		params = append(params, "name")
	}

	if environment.ApplicationID == 0 {
		params = append(params, "application_id")
	}

	if len(params) > 0 {
		return http.ErrorInvalidParam{Params: params}
	}

	return nil
}
