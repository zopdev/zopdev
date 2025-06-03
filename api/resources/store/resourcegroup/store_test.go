package resourcegroup

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"

	"github.com/zopdev/zopdev/api/resources/models"
)

func setup(t *testing.T) (*gofr.Context, *container.Mocks, *Store) {
	t.Helper()

	mockContainer, mocks := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}

	return ctx, mocks, New()
}

func TestStore_GetAllResourceGroups(t *testing.T) {
	ctx, mocks, store := setup(t)
	query := `SELECT id, name, description, cloud_account_id FROM resource_groups
                                               WHERE cloud_account_id = ? AND deleted_at IS NULL`
	rows := sqlmock.NewRows([]string{"id", "name", "description", "cloud_account_id"}).
		AddRow(1, "group1", "desc1", 123).
		AddRow(2, "group2", "desc2", 123)

	mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(123).WillReturnRows(rows)

	result, err := store.GetAllResourceGroups(ctx, 123)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "group1", result[0].Name)
}

func TestStore_GetAllResourceGroups_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	query := `SELECT id, name, description, cloud_account_id FROM resource_groups
                                               WHERE cloud_account_id = ? AND deleted_at IS NULL`
	mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(123).WillReturnError(assert.AnError)

	result, err := store.GetAllResourceGroups(ctx, 123)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestStore_GetResourceGroupByID(t *testing.T) {
	ctx, mocks, store := setup(t)
	query := `SELECT id, name, description, cloud_account_id FROM resource_groups
                                               WHERE id = ? AND cloud_account_id = ? AND deleted_at IS NULL`
	row := sqlmock.NewRows([]string{"id", "name", "description", "cloud_account_id"}).
		AddRow(1, "group1", "desc1", 123)
	mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(1, 123).WillReturnRows(row)

	result, err := store.GetResourceGroupByID(ctx, 123, 1)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
}

func TestStore_GetResourceGroupByID_NotFound(t *testing.T) {
	ctx, mocks, store := setup(t)
	query := `SELECT id, name, description, cloud_account_id FROM resource_groups
                                               WHERE id = ? AND cloud_account_id = ? AND deleted_at IS NULL`
	mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(2, 123).WillReturnError(sql.ErrNoRows)

	result, err := store.GetResourceGroupByID(ctx, 123, 2)

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestStore_CreateResourceGroup(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO resource_groups (name, description, cloud_account_id) VALUES (?, ?, ?)`).
		WithArgs("group1", "desc1", int64(123)).
		WillReturnResult(sqlmock.NewResult(10, 1))

	id, err := store.CreateResourceGroup(ctx, &models.RGCreate{Name: "group1", Description: "desc1", CloudAccountID: 123})

	require.NoError(t, err)
	assert.Equal(t, int64(10), id)
}

func TestStore_CreateResourceGroup_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO resource_groups (name, description, cloud_account_id) VALUES (?, ?, ?)`).
		WithArgs("group1", "desc1", int64(123)).
		WillReturnError(assert.AnError)

	id, err := store.CreateResourceGroup(ctx, &models.RGCreate{Name: "group1", Description: "desc1", CloudAccountID: 123})

	require.Error(t, err)
	assert.Equal(t, int64(0), id)
}

func TestStore_UpdateResourceGroup(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`UPDATE resource_groups SET name = ?, description = ? WHERE id = ?`).
		WithArgs("group1", "desc1", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.UpdateResourceGroup(ctx, &models.RGUpdate{ID: 1, Name: "group1", Description: "desc1"})

	require.NoError(t, err)
}

func TestStore_UpdateResourceGroup_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`UPDATE resource_groups SET name = ?, description = ? WHERE id = ?`).
		WithArgs("group1", "desc1", int64(1)).
		WillReturnError(assert.AnError)

	err := store.UpdateResourceGroup(ctx, &models.RGUpdate{ID: 1, Name: "group1", Description: "desc1"})

	require.Error(t, err)
}

func TestStore_DeleteResourceGroup(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`UPDATE resource_groups SET deleted_at = ? WHERE id = ?`).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mocks.SQL.Sqlmock.ExpectExec(`DELETE FROM resource_group_memberships WHERE group_id = ?`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.DeleteResourceGroup(ctx, 1)

	require.NoError(t, err)
}

func TestStore_DeleteResourceGroup_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`UPDATE resource_groups SET deleted_at = ? WHERE id = ?`).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnError(assert.AnError)

	err := store.DeleteResourceGroup(ctx, 1)

	require.Error(t, err)
}

func TestStore_GetResourceIDs(t *testing.T) {
	ctx, mocks, store := setup(t)
	rows := sqlmock.NewRows([]string{"resource_id"}).AddRow(1).AddRow(2)
	mocks.SQL.Sqlmock.ExpectQuery(`SELECT resource_id FROM resource_group_memberships WHERE group_id = ? ORDER BY  resource_id`).
		WithArgs(int64(1)).WillReturnRows(rows)

	ids, err := store.GetResourceIDs(ctx, 1)
	
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, ids)
}

func TestStore_GetResourceIDs_Error(t *testing.T) {
	ctx, mocks, store := setup(t)

	mocks.SQL.Sqlmock.ExpectQuery(`SELECT resource_id FROM resource_group_memberships WHERE group_id = ? ORDER BY  resource_id`).
		WithArgs(int64(1)).WillReturnError(assert.AnError)

	ids, err := store.GetResourceIDs(ctx, 1)

	require.Error(t, err)
	assert.Nil(t, ids)
}

func TestStore_AddResourcesToGroup(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO resource_group_memberships (resource_id, group_id) VALUES (?, ?)`).
		WithArgs(int64(1), int64(2)).WillReturnResult(sqlmock.NewResult(1, 1))
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO resource_group_memberships (resource_id, group_id) VALUES (?, ?)`).
		WithArgs(int64(2), int64(2)).WillReturnResult(sqlmock.NewResult(2, 1))

	err := store.AddResourcesToGroup(ctx, 2, []int64{1, 2})

	require.NoError(t, err)
}

func TestStore_AddResourcesToGroup_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO resource_group_memberships (resource_id, group_id) VALUES (?, ?)`).
		WithArgs(int64(1), int64(2)).WillReturnError(assert.AnError)

	err := store.AddResourcesToGroup(ctx, 2, []int64{1})

	require.Error(t, err)
}

func TestStore_RemoveResourceFromGroup(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`DELETE from resource_group_memberships WHERE group_id = ? AND resource_id = ?`).
		WithArgs(int64(2), int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.RemoveResourceFromGroup(ctx, 2, 1)

	require.NoError(t, err)
}

func TestStore_RemoveResourceFromGroup_Error(t *testing.T) {
	ctx, mocks, store := setup(t)
	mocks.SQL.Sqlmock.ExpectExec(`DELETE from resource_group_memberships WHERE group_id = ? AND resource_id = ?`).
		WithArgs(int64(2), int64(1)).WillReturnError(assert.AnError)

	err := store.RemoveResourceFromGroup(ctx, 2, 1)

	require.Error(t, err)
}
