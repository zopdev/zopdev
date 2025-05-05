package service

import (
	"github.com/zopdev/zopdev/api/environments/store"
	"gofr.dev/pkg/gofr"
)

// EnvironmentService defines the business logic layer for managing environments.
// It provides methods to fetch, add, and update environments.
type EnvironmentService interface {

	// FetchAll retrieves all environments associated with a specific application.
	//
	// Parameters:
	//   - ctx: The request context, which includes logging and request-scoped values.
	//   - applicationID: The unique identifier of the application whose environments need to be fetched.
	//
	// Returns:
	//   - []store.Environment: A slice of environments linked to the application.
	//   - error: An error if the operation fails.
	FetchAll(ctx *gofr.Context, applicationID int) ([]store.Environment, error)

	// Add creates a new environment and stores it in the datastore.
	//
	// Parameters:
	//   - ctx: The request context, which includes logging and request-scoped values.
	//   - environment: A pointer to the Environment struct containing the details to be added.
	//
	// Returns:
	//   - *store.Environment: A pointer to the newly created environment record, including its ID.
	//   - error: An error if the operation fails.
	Add(ctx *gofr.Context, environment *store.Environment) (*store.Environment, error)

	// Update modifies multiple environment records in the datastore.
	//
	// Parameters:
	//   - ctx: The request context, which includes logging and request-scoped values.
	//   - environments: A slice of Environment structs containing the updated details.
	//
	// Returns:
	//   - []store.Environment: A slice of the updated environment records.
	//   - error: An error if the operation fails.
	Update(ctx *gofr.Context, environments []store.Environment) ([]store.Environment, error)
}
