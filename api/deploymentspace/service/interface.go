package service

import (
	"gofr.dev/pkg/gofr"
)

// DeploymentSpaceService is an interface that provides methods for managing deployment spaces.
// It includes methods to add a deployment space and fetch deployment space details by environment ID.
type DeploymentSpaceService interface {
	// Add adds a new deployment space to the system.
	// It accepts the context, a DeploymentSpace object containing the deployment details, and an environment ID.
	// The method returns the created DeploymentSpace and any error encountered during the operation.
	//
	// Parameters:
	//   ctx - The GoFR context that carries request-specific data like SQL connections.
	//   deploymentSpace - The deployment space object to be added.
	//   environmentID - The ID of the environment in which the deployment space will be created.
	//
	// Returns:
	//   *DeploymentSpace - The newly added deployment space with updated details (including ID).
	//   error - Any error that occurs during the add operation.
	Add(ctx *gofr.Context, deploymentSpace *DeploymentSpace, environmentID int) (*DeploymentSpace, error)

	// Fetch fetches the deployment space details for a given environment ID.
	// It returns a DeploymentSpaceResp object, which includes the deployment space and the associated cluster.
	//
	// Parameters:
	//   ctx - The GoFR context that carries request-specific data like SQL connections.
	//   environmentID - The ID of the environment for which the deployment space details are to be fetched.
	//
	// Returns:
	//   *DeploymentSpaceResp - The deployment space response object that includes the deployment space and cluster.
	//   error - Any error encountered during the fetch operation.
	Fetch(ctx *gofr.Context, environmentID int) (*DeploymentSpaceResp, error)
	GetServices(ctx *gofr.Context, environmentID int) (any, error)
	GetDeployments(ctx *gofr.Context, environmentID int) (any, error)
	GetPods(ctx *gofr.Context, environmentID int) (any, error)
	GetCronJobs(ctx *gofr.Context, environmentID int) (any, error)
	GetServiceByName(ctx *gofr.Context, envID int, serviceName string) (any, error)
	GetDeploymentByName(ctx *gofr.Context, envID int, deploymentName string) (any, error)
	GetPodByName(ctx *gofr.Context, environmentID int, deploymentName string) (any, error)
	GetCronJobByName(ctx *gofr.Context, environmentID int, deploymentName string) (any, error)
}
