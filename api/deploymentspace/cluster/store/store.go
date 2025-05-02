package store

import "gofr.dev/pkg/gofr"

// Store implements the ClusterStore interface for managing cluster data in the storage system.
type Store struct{}

// New creates and returns a new instance of Store that implements ClusterStore.
func New() ClusterStore {
	return &Store{}
}

// Insert inserts a new cluster into the storage and returns the inserted cluster with its generated ID.
func (*Store) Insert(ctx *gofr.Context, cluster *Cluster) (*Cluster, error) {
	res, err := ctx.SQL.ExecContext(ctx, INSERTQUERY, cluster.DeploymentSpaceID, cluster.Identifier,
		cluster.Name, cluster.Region, cluster.ProviderID, cluster.Provider, cluster.Namespace.Name)
	if err != nil {
		return nil, err
	}

	cluster.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// GetByDeploymentSpaceID retrieves a cluster by its associated deployment space ID.
func (*Store) GetByDeploymentSpaceID(ctx *gofr.Context, deploymentSpaceID int) (*Cluster, error) {
	cluster := Cluster{}

	err := ctx.SQL.QueryRowContext(ctx, GETQUERY, deploymentSpaceID).Scan(&cluster.ID, &cluster.DeploymentSpaceID, &cluster.Identifier,
		&cluster.Name, &cluster.Region, &cluster.ProviderID, &cluster.Provider, &cluster.Namespace.Name, &cluster.CreatedAt, &cluster.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

// GetByCluster retrieves a cluster by its associated deployment space ID.
func (*Store) GetByCluster(ctx *gofr.Context, cluster *Cluster) (*Cluster, error) {
	err := ctx.SQL.QueryRowContext(ctx, GETBYCLUSTER, cluster.Provider, cluster.Name,
		cluster.Region, cluster.ProviderID, cluster.Namespace.Name).
		Scan(&cluster.ID, &cluster.DeploymentSpaceID, &cluster.Identifier, &cluster.Name,
			&cluster.Region, &cluster.ProviderID, &cluster.Provider, &cluster.Namespace.Name,
			&cluster.CreatedAt, &cluster.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}
