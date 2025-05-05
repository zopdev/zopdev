/*
Package deploymentspace provides an interface for interacting with deployment spaces.
It defines methods for fetching deployment space details by ID and adding new resources to the deployment space.
*/
package deploymentspace

import (
	"gofr.dev/pkg/gofr"
)

// DeploymentEntity defines the interface for managing deployment spaces.
// It provides methods to fetch details of a deployment space by its ID and add resources to the deployment space.
type DeploymentEntity interface {
	// FetchByDeploymentSpaceID retrieves a deployment space by its unique identifier.
	// It returns the resource details as an interface{}, or an error if the operation fails.
	//
	// Parameters:
	//   ctx - The context of the request.
	//   id  - The unique identifier of the deployment space.
	//
	// Returns:
	//   An interface{} containing the deployment space resource details, or an error if fetching the resource fails.
	FetchByDeploymentSpaceID(ctx *gofr.Context, id int) (interface{}, error)

	// Add adds a new resource to the deployment space.
	// It returns the newly added resource details as an interface{}, or an error if the operation fails.
	//
	// Parameters:
	//   ctx      - The context of the request.
	//   resource - The resource to be added to the deployment space.
	//
	// Returns:
	//   An interface{} containing the newly added resource details, or an error if the addition fails.
	Add(ctx *gofr.Context, resource any) (interface{}, error)

	// DuplicateCheck check if deploymentspace is already configure with any environment
	// Parameters:
	//   ctx      - The context of the request.
	//   resource - The resource to be added to the deployment space.
	//
	// Returns:
	//  An interface{} containing the resource details, or an error if the duplicate is found
	DuplicateCheck(ctx *gofr.Context, data any) (interface{}, error)
}
