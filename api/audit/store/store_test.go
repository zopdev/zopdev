package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestStore_GetLastRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)

	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	mockCloudAccountID := int64(1)
	mockRule := "test_rule"
	mockResult := &Result{
		ID:             1,
		CloudAccountID: mockCloudAccountID,
		RuleID:         mockRule,
		Result:         &ResultData{Data: []Items{{"instance1", "passing", nil}}},
		EvaluatedAt:    time.Now(),
	}

	mocks.SQL.Sqlmock.ExpectQuery("SELECT id, cloud_account_id, rule_id, result, evaluated_at "+
		"FROM results WHERE cloud_account_id = ? AND rule_id = ? ORDER BY evaluated_at DESC LIMIT 1").
		WithArgs(mockResult.CloudAccountID, mockResult.RuleID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cloud_account_id", "rule_id", "result", "evaluated_at"}).
			AddRow(1, mockResult.CloudAccountID, mockResult.RuleID, mockResult.Result, mockResult.EvaluatedAt))

	res, err := store.GetLastRun(ctx, mockCloudAccountID, mockRule)
	require.NoError(t, err)
	assert.Equal(t, res, mockResult)

	// no rows
	mocks.SQL.Sqlmock.ExpectQuery("SELECT id, cloud_account_id, rule_id, result, evaluated_at "+
		"FROM results WHERE cloud_account_id = ? AND rule_id = ? ORDER BY evaluated_at DESC LIMIT 1").
		WithArgs(mockResult.CloudAccountID, mockResult.RuleID).WillReturnError(sql.ErrNoRows)

	res, err = store.GetLastRun(ctx, mockCloudAccountID, mockRule)
	require.NoError(t, err)
	assert.Nil(t, res)

	// error case
	mocks.SQL.Sqlmock.ExpectQuery("SELECT id, cloud_account_id, rule_id, result, evaluated_at "+
		"FROM results WHERE cloud_account_id = ? AND rule_id = ? ORDER BY evaluated_at DESC LIMIT 1").
		WithArgs(mockResult.CloudAccountID, mockResult.RuleID).WillReturnError(sql.ErrConnDone)

	mocks.Metrics.EXPECT().IncrementCounter(ctx, "db_error_count", "audit_store", "GetLastRun", "error", sql.ErrConnDone.Error())

	res, err = store.GetLastRun(ctx, mockCloudAccountID, mockRule)
	require.Error(t, err)
	assert.Nil(t, res)
}

func TestStore_CreatePending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)

	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	mockResult := &Result{
		CloudAccountID: 1,
		RuleID:         "test_rule",
		EvaluatedAt:    time.Now(),
	}

	// Mock successful insert
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO results (cloud_account_id, rule_id, evaluated_at) VALUES (?, ?, ?)`).
		WithArgs(mockResult.CloudAccountID, mockResult.RuleID, mockResult.EvaluatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := store.CreatePending(ctx, mockResult)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.ID)

	// Mock insert error
	mocks.SQL.Sqlmock.ExpectExec(`INSERT INTO results (cloud_account_id, rule_id, evaluated_at) VALUES (?, ?, ?)`).
		WithArgs(mockResult.CloudAccountID, mockResult.RuleID, mockResult.EvaluatedAt).
		WillReturnError(sql.ErrConnDone)
	mocks.Metrics.EXPECT().IncrementCounter(ctx, "db_error_count", "audit_store", "CreatePending", "error", sql.ErrConnDone.Error())

	res, err = store.CreatePending(ctx, mockResult)
	require.Error(t, err)
	assert.Nil(t, res)
}

func TestStore_UpdateResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)

	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	mockResult := &Result{
		ID:     1,
		Result: &ResultData{Data: []Items{{"instance1", "passing", nil}}},
	}

	// Mock successful update
	mocks.SQL.Sqlmock.ExpectExec("UPDATE results SET result = ? WHERE id = ?").
		WithArgs(mockResult.Result, mockResult.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := store.UpdateResult(ctx, mockResult)
	require.NoError(t, err)

	// Mock update error
	mocks.SQL.Sqlmock.ExpectExec("UPDATE results SET result = ? WHERE id = ?").
		WithArgs(mockResult.Result, mockResult.ID).
		WillReturnError(sql.ErrConnDone)
	mocks.Metrics.EXPECT().IncrementCounter(ctx, "db_error_count", "audit_store", "UpdateResult", "error", sql.ErrConnDone.Error())

	err = store.UpdateResult(ctx, mockResult)
	require.Error(t, err)
}
