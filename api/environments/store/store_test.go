package store

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestStore_Insert(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	store := &Store{}

	environment := &Environment{
		Name:          "Test Environment",
		Level:         1,
		ApplicationID: 1,
	}

	mock.SQL.ExpectExec(INSERTQUERY).
		WithArgs(environment.Name, environment.Level, environment.ApplicationID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := store.Insert(ctx, environment)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(1), res.ID)

	require.NoError(t, mock.SQL.ExpectationsWereMet())
}

func TestStore_GetALL(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	store := &Store{}

	applicationID := 1
	environments := []Environment{
		{
			ID:            1,
			Name:          "Test Environment 1",
			Level:         2,
			ApplicationID: int64(applicationID),
			CreatedAt:     time.Now().String(),
			UpdatedAt:     time.Now().String(),
		},
		{
			ID:            2,
			Name:          "Test Environment 2",
			Level:         1,
			ApplicationID: int64(applicationID),
			CreatedAt:     time.Now().String(),
			UpdatedAt:     time.Now().String(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "level", "application_id", "created_at", "updated_at"}).
		AddRow(environments[0].ID, environments[0].Name, environments[0].Level,
			environments[0].ApplicationID, environments[0].CreatedAt, environments[0].UpdatedAt).
		AddRow(environments[1].ID, environments[1].Name, environments[1].Level,
			environments[1].ApplicationID, environments[1].CreatedAt, environments[1].UpdatedAt)

	mock.SQL.ExpectQuery(GETALLQUERY).WithArgs(applicationID).WillReturnRows(rows)

	res, err := store.GetALL(ctx, applicationID)
	require.NoError(t, err)
	require.Len(t, res, 2)

	require.NoError(t, mock.SQL.ExpectationsWereMet())
}

func TestStore_GetByName(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	store := &Store{}

	applicationID := 1
	name := "Test Environment"
	environment := &Environment{
		ID:            1,
		Name:          name,
		Level:         1,
		ApplicationID: int64(applicationID),
		CreatedAt:     time.Now().String(),
		UpdatedAt:     time.Now().String(),
	}

	row := sqlmock.NewRows([]string{"id", "name", "level", "application_id", "created_at", "updated_at"}).
		AddRow(environment.ID, environment.Name, environment.Level, environment.ApplicationID, environment.CreatedAt, environment.UpdatedAt)

	mock.SQL.ExpectQuery(GETBYNAMEQUERY).WithArgs(name, applicationID).WillReturnRows(row)

	res, err := store.GetByName(ctx, applicationID, name)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, name, res.Name)

	require.NoError(t, mock.SQL.ExpectationsWereMet())
}

func TestStore_Update(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	store := &Store{}

	environment := &Environment{
		ID:            1,
		Name:          "Updated Environment",
		Level:         2,
		ApplicationID: 1,
	}

	mock.SQL.ExpectExec(UPDATEQUERY).
		WithArgs(environment.Name, environment.Level, environment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := store.Update(ctx, environment)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "Updated Environment", res.Name)

	require.NoError(t, mock.SQL.ExpectationsWereMet())
}
