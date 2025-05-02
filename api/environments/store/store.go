package store

import (
	"time"

	"gofr.dev/pkg/gofr"
)

type Store struct {
}

// New creates and returns a new instance of EnvironmentStore.
//
// Returns:
//   - EnvironmentStore: A new instance of the store implementation.
func New() EnvironmentStore {
	return &Store{}
}

// Insert inserts a new environment record into the datastore.
//
// Parameters:
//   - ctx: The request context, which includes database connection and logging.
//   - environment: A pointer to the Environment struct containing the details to be inserted.
//
// Returns:
//   - *Environment: A pointer to the newly inserted environment record, including the generated ID.
//   - error: An error if the operation fails.
func (*Store) Insert(ctx *gofr.Context, environment *Environment) (*Environment, error) {
	res, err := ctx.SQL.ExecContext(ctx, INSERTQUERY, environment.Name, environment.Level, environment.ApplicationID)
	if err != nil {
		return nil, err
	}

	environment.ID, _ = res.LastInsertId()

	environment.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	return environment, nil
}

// GetALL retrieves all environments associated with a specific application.
//
// Parameters:
//   - ctx: The request context, which includes database connection and logging.
//   - applicationID: The unique identifier of the application whose environments need to be fetched.
//
// Returns:
//   - []Environment: A slice of environments associated with the specified application.
//   - error: An error if the operation fails or the query returns no results.
func (*Store) GetALL(ctx *gofr.Context, applicationID int) ([]Environment, error) {
	environments := []Environment{}

	rows, err := ctx.SQL.QueryContext(ctx, GETALLQUERY, applicationID)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var environment Environment

		err := rows.Scan(&environment.ID, &environment.Name, &environment.Level, &environment.ApplicationID,
			&environment.CreatedAt, &environment.UpdatedAt)
		if err != nil {
			return nil, err
		}

		environments = append(environments, environment)
	}

	return environments, nil
}

// GetByName retrieves a specific environment by its name for a given application.
//
// Parameters:
//   - ctx: The request context, which includes database connection and logging.
//   - applicationID: The unique identifier of the application to which the environment belongs.
//   - name: The name of the environment to retrieve.
//
// Returns:
//   - *Environment: A pointer to the environment matching the specified name.
//   - error: An error if the operation fails or no environment is found.
func (*Store) GetByName(ctx *gofr.Context, applicationID int, name string) (*Environment, error) {
	row := ctx.SQL.QueryRowContext(ctx, GETBYNAMEQUERY, name, applicationID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var environment Environment

	err := row.Scan(&environment.ID, &environment.Name, &environment.Level, &environment.ApplicationID,
		&environment.CreatedAt, &environment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

// Update modifies an existing environment record in the datastore.
//
// Parameters:
//   - ctx: The request context, which includes database connection and logging.
//   - environment: A pointer to the Environment struct containing the updated details.
//
// Returns:
//   - *Environment: A pointer to the updated environment record.
//   - error: An error if the operation fails.
func (*Store) Update(ctx *gofr.Context, environment *Environment) (*Environment, error) {
	_, err := ctx.SQL.ExecContext(ctx, UPDATEQUERY, environment.Name, environment.Level, environment.ID)
	if err != nil {
		return nil, err
	}

	return environment, nil
}

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
func (*Store) GetMaxLevel(ctx *gofr.Context, applicationID int) (int, error) {
	var maxLevel int

	err := ctx.SQL.QueryRowContext(ctx, GETMAXLEVEL, applicationID).Scan(&maxLevel)
	if err != nil {
		return 0, err
	}

	return maxLevel, nil
}
