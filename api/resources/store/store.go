package store

import "gofr.dev/pkg/gofr"

type Store struct{}

func New() *Store { return &Store{} }

func (s *Store) InsertResource(ctx *gofr.Context, res Resource) error {
	_, err := ctx.SQL.ExecContext(ctx,
		`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type) VALUES (?, ?, ?, ?, ?, ?)`,
		res.UID, res.Name, res.State, res.CloudAccountID, res.CloudProvider, res.Type)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetResources(ctx *gofr.Context, cloudAccountID int64, resourceType []string) ([]Resource, error) {
	var (
		resources []Resource
		args      = make([]any, 0, 4)
	)

	// If the resources are given as a parameter, we can use that to filter.
	// Form the IN clause, otherwise we fetch all resources for the given cloud account ID.
	inClause := ``
	args = append(args, cloudAccountID)

	if len(resourceType) > 0 {
		inClause = ` AND resource_type IN (`

		for _, res := range resourceType {
			if inClause != `` {
				inClause += `, `
			}

			inClause += `?`
			args = append(args, res)
		}

		inClause += `)`
	}

	rows, err := ctx.SQL.QueryContext(ctx, `SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at FROM resources WHERE cloud_account_id = ?`+inClause+`ORDER BY resource_uid`, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var res Resource
		if er := rows.Scan(&res.ID, &res.UID, &res.Name, &res.State,
			&res.CloudAccountID, &res.CloudProvider, &res.Type, &res.CreatedAt, &res.UpdatedAt); er != nil {
			return nil, er
		}

		resources = append(resources, res)
	}

	return resources, nil
}

func (s *Store) UpdateResource(ctx *gofr.Context, res Resource) error {
	_, err := ctx.SQL.ExecContext(ctx, `UPDATE resources SET resource_type = ?, state = ? WHERE id = ?`,
		res.Type, res.State, res.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) RemoveResource(ctx *gofr.Context, id int64) error {
	_, err := ctx.SQL.ExecContext(ctx, `DELETE FROM resources WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}
