/*
Package store provides an interface for interacting with deployment spaces in a store.
It defines methods for inserting a new deployment space and fetching a deployment space by its associated environment ID.
*/
package store

import "gofr.dev/pkg/gofr"

// DeploymentSpaceStore defines the interface for managing deployment spaces in a data store.
// It provides methods for inserting new deployment spaces and retrieving deployment spaces by environment ID.
type DeploymentSpaceStore interface {
	// Insert adds a new deployment space to the store.
	// It returns the inserted deployment space or an error if the operation fails.
	//
	// Parameters:
	//   ctx             - The context of the request.
	//   deploymentSpace - The deployment space to be inserted.
	//
	// Returns:
	//   The inserted deployment space or an error if the operation fails.
	Insert(ctx *gofr.Context, deploymentSpace *DeploymentSpace) (*DeploymentSpace, error)

	// GetByEnvironmentID retrieves a deployment space associated with the given environment ID.
	// It returns the deployment space or an error if the operation fails.
	//
	// Parameters:
	//   ctx          - The context of the request.
	//   environmentID - The unique identifier of the environment.
	//
	// Returns:
	//   The deployment space associated with the environment ID or an error if fetching the resource fails.
	GetByEnvironmentID(ctx *gofr.Context, environmentID int) (*DeploymentSpace, error)
}
