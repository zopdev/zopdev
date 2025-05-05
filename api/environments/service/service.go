package service

import (
	"database/sql"
	"errors"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/deploymentspace/service"
	"github.com/zopdev/zopdev/api/environments/store"
)

type Service struct {
	store                  store.EnvironmentStore
	deploymentSpaceService service.DeploymentSpaceService
}

// New creates a new instance of the EnvironmentService.
//
// Parameters:
//   - enStore: The EnvironmentStore implementation for managing datastore operations.
//   - deploySvc: The DeploymentSpaceService implementation for handling deployment space logic.
//
// Returns:
//   - EnvironmentService: A new instance of the service layer.
func New(enStore store.EnvironmentStore, deploySvc service.DeploymentSpaceService) EnvironmentService {
	return &Service{store: enStore, deploymentSpaceService: deploySvc}
}

// Add adds a new environment to the datastore.
//
// The method first checks if an environment with the same name already exists for the application.
// If no such environment exists, it creates a new record.
//
// Parameters:
//   - ctx: The request context, which includes logging and request-scoped values.
//   - environment: A pointer to the Environment struct containing the details to be added.
//
// Returns:
//   - *store.Environment: A pointer to the newly created environment record, including its ID.
//   - error:
//   - An error if the operation fails.
//   - http.ErrorEntityAlreadyExist if the environment already exists.
func (s *Service) Add(ctx *gofr.Context, environment *store.Environment) (*store.Environment, error) {
	tempEnvironment, err := s.store.GetByName(ctx, int(environment.ApplicationID), environment.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			return nil, err
		}
	}

	if tempEnvironment != nil {
		return nil, http.ErrorEntityAlreadyExist{}
	}

	maxLevel, err := s.store.GetMaxLevel(ctx, int(environment.ApplicationID))
	if err != nil {
		return nil, err
	}

	environment.Level = maxLevel + 1

	return s.store.Insert(ctx, environment)
}

// FetchAll retrieves all environments for a specific application and populates their deployment spaces.
//
// The method fetches environments from the datastore and augments them with deployment space details
// retrieved from the DeploymentSpaceService.
//
// Parameters:
//   - ctx: The request context, which includes logging and request-scoped values.
//   - applicationID: The unique identifier of the application whose environments need to be fetched.
//
// Returns:
//   - []store.Environment: A slice of environments, each augmented with deployment space details.
//   - error: An error if the operation fails.

func (s *Service) FetchAll(ctx *gofr.Context, applicationID int) ([]store.Environment, error) {
	environments, err := s.store.GetALL(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	for i := range environments {
		deploymentSpace, err := s.deploymentSpaceService.Fetch(ctx, int(environments[i].ID))
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			return nil, err
		}

		if deploymentSpace != nil {
			deploymentSpaceResp := DeploymentSpaceResponse{
				Name: deploymentSpace.DeploymentSpace.CloudAccountName,
				Next: &DeploymentSpaceResponse{
					Name: deploymentSpace.DeploymentSpace.Type,
					Next: &DeploymentSpaceResponse{
						Name: deploymentSpace.Cluster.Name,
						Next: &DeploymentSpaceResponse{
							Name: deploymentSpace.Cluster.Namespace.Name,
						},
					},
				},
			}
			environments[i].DeploymentSpace = deploymentSpaceResp
		}
	}

	return environments, nil
}

// Update modifies existing environment records in the datastore.
//
// The method iterates through the list of environments, updates each one, and returns the updated records.
//
// Parameters:
//   - ctx: The request context, which includes logging and request-scoped values.
//   - environments: A slice of Environment structs containing the updated details.
//
// Returns:
//   - []store.Environment: A slice of the updated environment records.
//   - error: An error if the operation fails.
func (s *Service) Update(ctx *gofr.Context, environments []store.Environment) ([]store.Environment, error) {
	for i := range environments {
		env, err := s.store.Update(ctx, &environments[i])
		if err != nil {
			return nil, err
		}

		environments[i] = *env
	}

	return environments, nil
}
