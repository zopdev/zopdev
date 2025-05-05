package store

import "gofr.dev/pkg/gofr"

// EnvironmentStore provides an abstraction for managing environment-related operations
// in the datastore. It defines methods to insert, fetch, and update environments.
type EnvironmentStore interface {

	// Insert adds a new environment to the datastore.
	//
	// Parameters:
	//   - ctx: The context for the request, which includes logging and request-scoped values.
	//   - environment: A pointer to the Environment struct containing the details of the environment to insert.
	//
	// Returns:
	//   - *Environment: A pointer to the inserted environment record.
	//   - error: An error if the operation fails.
	Insert(ctx *gofr.Context, environment *Environment) (*Environment, error)

	// GetALL fetches all environments for a specific application from the datastore.
	//
	// Parameters:
	//   - ctx: The context for the request, which includes logging and request-scoped values.
	//   - applicationID: The unique identifier of the application whose environments need to be fetched.
	//
	// Returns:
	//   - []Environment: A slice of environments associated with the application.
	//   - error: An error if the operation fails.
	GetALL(ctx *gofr.Context, applicationID int) ([]Environment, error)

	// GetByName fetches a specific environment by its name for a given application.
	//
	// Parameters:
	//   - ctx: The context for the request, which includes logging and request-scoped values.
	//   - applicationID: The unique identifier of the application to which the environment belongs.
	//   - name: The name of the environment to fetch.
	//
	// Returns:
	//   - *Environment: A pointer to the environment matching the given name.
	//   - error: An error if the operation fails or no environment is found.
	GetByName(ctx *gofr.Context, applicationID int, name string) (*Environment, error)

	// Update updates the details of an existing environment in the datastore.
	//
	// Parameters:
	//   - ctx: The context for the request, which includes logging and request-scoped values.
	//   - environment: A pointer to the Environment struct containing the updated details.
	//
	// Returns:
	//   - *Environment: A pointer to the updated environment record.
	//   - error: An error if the operation fails.
	Update(ctx *gofr.Context, environment *Environment) (*Environment, error)

	// GetMaxLevel retrieves the maximum environment `level` for a specified application from the `environment` table.
	// It considers only active records where `deleted_at IS NULL`.
	//
	// Parameters:
	//   - ctx (*gofr.Context): Context for request handling, logging, and tracing.
	//   - applicationID (int): The unique identifier for the application.
	//
	// Returns:
	//   - int: The highest environment level associated with the given application ID.
	//   - error: An error if the query fails or no matching data exists.
	GetMaxLevel(ctx *gofr.Context, applicationID int) (int, error)
}
