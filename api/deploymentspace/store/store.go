/*
Package store provides an implementation of the DeploymentSpaceStore interface for managing deployment spaces.
The Store struct implements methods to insert a new deployment space and retrieve a deployment space by its environment ID.
*/
package store

import "gofr.dev/pkg/gofr"

type Store struct{}

// New creates and returns a new instance of Store, which implements the DeploymentSpaceStore interface.
func New() DeploymentSpaceStore {
	return &Store{}
}

// Insert inserts a new deployment space into the data store.
// It takes a context and a pointer to the deployment space to be inserted.
// Upon successful insertion, it returns the newly inserted deployment space, including its generated ID.
// If there is an error during insertion, it returns nil and the error.
func (*Store) Insert(ctx *gofr.Context, deploymentSpace *DeploymentSpace) (*DeploymentSpace, error) {
	res, err := ctx.SQL.ExecContext(ctx, INSERTQUERY, deploymentSpace.CloudAccountID, deploymentSpace.EnvironmentID, deploymentSpace.Type)
	if err != nil {
		return nil, err
	}

	deploymentSpace.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return deploymentSpace, nil
}

// GetByEnvironmentID retrieves a deployment space based on the given environment ID.
// It queries the data store using the environment ID and populates the provided DeploymentSpace object with the retrieved data.
// If successful, it returns a pointer to the populated DeploymentSpace. If there is an error, it returns nil and the error.
func (*Store) GetByEnvironmentID(ctx *gofr.Context, environmentID int) (*DeploymentSpace, error) {
	deploymentSpace := DeploymentSpace{}

	err := ctx.SQL.QueryRowContext(ctx, GETQUERYBYENVID, environmentID).Scan(&deploymentSpace.ID,
		&deploymentSpace.CloudAccountID, &deploymentSpace.EnvironmentID, &deploymentSpace.Type,
		&deploymentSpace.CreatedAt, &deploymentSpace.UpdatedAt, &deploymentSpace.CloudAccountName, &deploymentSpace.EnvironmentName)
	if err != nil {
		return nil, err
	}

	return &deploymentSpace, nil
}
