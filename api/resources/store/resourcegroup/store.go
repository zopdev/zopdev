package resourcegroup

import (
	"database/sql"
	"errors"
	"time"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

// GetAllResourceGroups retrieves all resource groups from the database.
func (*Store) GetAllResourceGroups(ctx *gofr.Context, cloudAccID int64) ([]models.ResourceGroup, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		`SELECT id, name, description, cloud_account_id FROM resource_groups 
                                               WHERE cloud_account_id = ? AND deleted_at IS NULL`, cloudAccID)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	var resourceGroups []models.ResourceGroup

	for rows.Next() {
		var resourceGroup models.ResourceGroup

		if er := rows.Scan(&resourceGroup.ID, &resourceGroup.Name,
			&resourceGroup.Description, &resourceGroup.CloudAccountID); er != nil {
			return nil, er
		}

		resourceGroups = append(resourceGroups, resourceGroup)
	}

	return resourceGroups, nil
}

// GetResourceGroupByID retrieves a resource group by its ID from the database.
func (*Store) GetResourceGroupByID(ctx *gofr.Context, cloudAccID, id int64) (*models.ResourceGroup, error) {
	row := ctx.SQL.QueryRowContext(ctx,
		`SELECT id, name, description, cloud_account_id FROM resource_groups 
                                               WHERE id = ? AND cloud_account_id = ? AND deleted_at IS NULL`, id, cloudAccID)

	var resourceGroup models.ResourceGroup

	err := row.Scan(&resourceGroup.ID, &resourceGroup.Name,
		&resourceGroup.Description, &resourceGroup.CloudAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No resource group found
		}

		return nil, err // Other error
	}

	return &resourceGroup, nil
}

// CreateResourceGroup inserts a new resource group into the database and returns its ID.
func (*Store) CreateResourceGroup(ctx *gofr.Context, resourceGroup *models.RGCreate) (int64, error) {
	result, err := ctx.SQL.ExecContext(ctx,
		`INSERT INTO resource_groups (name, description, cloud_account_id) VALUES (?, ?, ?)`,
		resourceGroup.Name, resourceGroup.Description, resourceGroup.CloudAccountID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateResourceGroup updates an existing resource group in the database.
func (*Store) UpdateResourceGroup(ctx *gofr.Context, resourceGroup *models.RGUpdate) error {
	_, err := ctx.SQL.ExecContext(ctx, `UPDATE resource_groups SET name = ?, description = ? WHERE id = ?`,
		resourceGroup.Name, resourceGroup.Description, resourceGroup.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteResourceGroup deletes a resource group from the database by its ID.
func (*Store) DeleteResourceGroup(ctx *gofr.Context, id int64) error {
	_, err := ctx.SQL.ExecContext(ctx, `UPDATE resource_groups SET deleted_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

// GetResourceIDs retrieves all resource IDs associated with a given resource group ID.
func (*Store) GetResourceIDs(ctx *gofr.Context, id int64) ([]int64, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		`SELECT resource_id FROM resource_group_memberships WHERE group_id = ? ORDER BY  resource_id`, id)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	var resIDs []int64

	for rows.Next() {
		var ids int64

		er := rows.Scan(&ids)
		if er != nil {
			return nil, er
		}

		resIDs = append(resIDs, ids)
	}

	return resIDs, nil
}

// AddResourcesToGroup adds a resource to a resource group by inserting a record into the membership table.
func (*Store) AddResourcesToGroup(ctx *gofr.Context, groupID int64, resourceIDs []int64) error {
	for _, resourceID := range resourceIDs {
		_, err := ctx.SQL.ExecContext(ctx,
			`INSERT INTO resource_group_memberships (resource_id, group_id) VALUES (?, ?)`,
			resourceID, groupID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveResourceFromGroup removes a resource from a resource group by deleting the record from the membership table.
func (*Store) RemoveResourceFromGroup(ctx *gofr.Context, groupID, resourceID int64) error {
	_, err := ctx.SQL.ExecContext(ctx,
		`DELETE from resource_group_memberships WHERE group_id = ? AND resource_id = ?`,
		groupID, resourceID)
	if err != nil {
		return err
	}

	return nil
}
