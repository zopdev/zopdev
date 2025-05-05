// Package service provides the business logic for managing clusters in a deployment space.

package service

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/zopdev/zopdev/api/deploymentspace"
	"github.com/zopdev/zopdev/api/deploymentspace/cluster/store"
	"gofr.dev/pkg/gofr"
)

// Service implements the deploymentspace.DeploymentEntity interface to manage clusters.
type Service struct {
	store store.ClusterStore
}

// New initializes and returns a new Service with the provided ClusterStore.
func New(str store.ClusterStore) deploymentspace.DeploymentEntity {
	return &Service{
		store: str,
	}
}

var (
	// errNamespaceAlreadyInUse indicates that the namespace is already associated with another environment.
	errNamespaceAlreadyInUSe = errors.New("namespace already in use")
)

// FetchByDeploymentSpaceID retrieves a cluster by its deployment space ID.
//
// Parameters:
//
//	ctx - The GoFR context carrying request-specific data.
//	id - The ID of the deployment space whose cluster is being fetched.
//
// Returns:
//
//	interface{} - The retrieved cluster details.
//	error - Any error encountered during the fetch operation.
func (s *Service) FetchByDeploymentSpaceID(ctx *gofr.Context, id int) (interface{}, error) {
	cluster, err := s.store.GetByDeploymentSpaceID(ctx, id)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// Add adds a new cluster to the storage.
//
// Parameters:
//
//	ctx - The GoFR context carrying request-specific data.
//	data - The cluster details to be added in a serializable format.
//
// Returns:
//
//	interface{} - The added cluster details.
//	error - Any error encountered during the add operation.
func (s *Service) Add(ctx *gofr.Context, data any) (interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	cluster := store.Cluster{}

	err = json.Unmarshal(bytes, &cluster)
	if err != nil {
		return nil, err
	}

	resp, err := s.store.Insert(ctx, &cluster)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DuplicateCheck checks if a cluster with the same details already exists.
//
// Parameters:
//
//	ctx - The GoFR context carrying request-specific data.
//	data - The cluster details to check for duplication in a serializable format.
//
// Returns:
//
//	interface{} - nil if no duplicate is found, otherwise an error indicating the namespace is already in use.
//	error - Any other error encountered during the check operation.
func (s *Service) DuplicateCheck(ctx *gofr.Context, data any) (interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	cluster := store.Cluster{}

	err = json.Unmarshal(bytes, &cluster)
	if err != nil {
		return nil, err
	}

	resp, err := s.store.GetByCluster(ctx, &cluster)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if resp != nil {
		return nil, errNamespaceAlreadyInUSe
	}

	return nil, nil
}
