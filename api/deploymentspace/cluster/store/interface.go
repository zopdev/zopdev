package store

import "gofr.dev/pkg/gofr"

// ClusterStore defines methods for interacting with cluster data in the storage system.
type ClusterStore interface {
	// Insert inserts a new cluster into the storage.
	//
	// Parameters:
	//   ctx - The request context.
	//   cluster - The Cluster object to be inserted.
	//
	// Returns:
	//   *Cluster - The inserted cluster with generated fields.
	//   error - Any error encountered during insertion.
	Insert(ctx *gofr.Context, cluster *Cluster) (*Cluster, error)

	// GetByDeploymentSpaceID retrieves a cluster by its deploymentSpaceID.
	//
	// Parameters:
	//   ctx - The request context.
	//   deploymentSpaceID - The deployment space ID.
	//
	// Returns:
	//   *Cluster - The cluster associated with the ID, or nil if not found.
	//   error - Any error encountered during retrieval.
	GetByDeploymentSpaceID(ctx *gofr.Context, deploymentSpaceID int) (*Cluster, error)

	// GetByCluster retrieves a cluster by its cluster configs.
	//
	// Parameters:
	//   ctx - The request context.
	//   cluster - cluster details
	//
	// Returns:
	//   *Cluster - The cluster associated with the ID, or nil if not found.
	//   error - Any error encountered during retrieval.
	GetByCluster(ctx *gofr.Context, cluster *Cluster) (*Cluster, error)
}
