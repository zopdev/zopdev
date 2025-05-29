package resource

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

// maxResTypes is the maximum number of resource types that are supported.
const maxResTypes = 10

type Store struct{}

func New() *Store { return &Store{} }

func (*Store) InsertResource(ctx *gofr.Context, res *models.Instance) error {
	_, err := ctx.SQL.ExecContext(ctx,
		`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type) VALUES (?, ?, ?, ?, ?, ?)`,
		res.UID, res.Name, res.Status, res.CloudAccount.ID, res.CloudAccount.Type, res.Type)
	if err != nil {
		return err
	}

	return nil
}

// GetResources fetches resources for a given cloud account ID.
// IMP: The returned result is sorted by resource UID. This is to ensure that the resources are returned in a consistent order.
// The service layer can use this to compare the resources fetched from the cloud provider with the resources stored in the database.
func (*Store) GetResources(ctx *gofr.Context, cloudAccountID int64, resourceType []string) ([]models.Instance, error) {
	var (
		resources []models.Instance
		args      = make([]any, 0, maxResTypes)
	)

	// Form the IN clause, otherwise we fetch all resources for the given cloud account ID.
	inClause := ``

	args = append(args, cloudAccountID)

	if len(resourceType) > 0 {
		inClause = ` AND resource_type IN (`

		for _, res := range resourceType {
			inClause += `?, `

			args = append(args, res)
		}

		inClause = inClause[:len(inClause)-2] // Remove the last comma

		inClause += `)`
	}

	rows, err := ctx.SQL.QueryContext(ctx, `SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at 
		FROM resources WHERE cloud_account_id = ?`+inClause+` ORDER BY resource_uid`, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var res models.Instance
		if er := rows.Scan(&res.ID, &res.UID, &res.Name, &res.Status,
			&res.CloudAccount.ID, &res.CloudAccount.Type, &res.Type, &res.CreatedAt, &res.UpdatedAt); er != nil {
			return nil, er
		}

		resources = append(resources, res)
	}

	return resources, nil
}

func (*Store) UpdateResource(ctx *gofr.Context, res *models.Instance) error {
	_, err := ctx.SQL.ExecContext(ctx, `UPDATE resources SET state = ? WHERE id = ?`,
		res.Status, res.ID)
	if err != nil {
		return err
	}

	return nil
}

func (*Store) RemoveResource(ctx *gofr.Context, id int64) error {
	_, err := ctx.SQL.ExecContext(ctx, `DELETE FROM resources WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}
