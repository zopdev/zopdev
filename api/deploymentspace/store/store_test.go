package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestInsertDeploymentSpace(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	deploymentSpace := &DeploymentSpace{
		CloudAccountID: 1,
		EnvironmentID:  1,
		Type:           "test-type",
	}

	testCases := []struct {
		name            string
		deploymentSpace *DeploymentSpace
		expectedError   bool
		mockBehavior    func()
	}{
		{
			name:            "success",
			deploymentSpace: deploymentSpace,
			expectedError:   false,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(deploymentSpace.CloudAccountID, deploymentSpace.EnvironmentID, deploymentSpace.Type).
					WillReturnResult(sqlmock.NewResult(123, 1))
			},
		},
		{
			name:            "error inserting deployment space",
			deploymentSpace: deploymentSpace,
			expectedError:   true,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(deploymentSpace.CloudAccountID, deploymentSpace.EnvironmentID, deploymentSpace.Type).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			result, err := store.Insert(ctx, tc.deploymentSpace)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, int64(123), result.ID)
			}
		})
	}
}

func TestGetDeploymentSpaceByEnvID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	expectedDeploymentSpace := &DeploymentSpace{
		ID:               1,
		CloudAccountID:   1,
		EnvironmentID:    1,
		Type:             "test-type",
		CloudAccountName: "Test Cloud Account",
		CreatedAt:        "2023-12-11T00:00:00Z",
		UpdatedAt:        "2023-12-11T00:00:00Z",
		EnvironmentName:  "hello",
	}

	testCases := []struct {
		name          string
		environmentID int
		mockBehavior  func()
		expectedError bool
		expectedSpace *DeploymentSpace
	}{
		{
			name:          "success",
			environmentID: 1,
			expectedError: false,
			mockBehavior: func() {
				mockRow := sqlmock.NewRows([]string{"id", "cloud_account_id", "environment_id", "type", "created_at", "updated_at",
					"cloud_account_name", "ev_name"}).
					AddRow(1, 1, 1, "test-type", "2023-12-11T00:00:00Z", "2023-12-11T00:00:00Z", "Test Cloud Account", "hello")
				mock.SQL.ExpectQuery(GETQUERYBYENVID).
					WithArgs(1).
					WillReturnRows(mockRow)
			},
			expectedSpace: expectedDeploymentSpace,
		},
		{
			name:          "no deployment space found",
			environmentID: 1,
			expectedError: true,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETQUERYBYENVID).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedSpace: nil,
		},
		{
			name:          "error on query execution",
			environmentID: 1,
			expectedError: true,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETQUERYBYENVID).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedSpace: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			result, err := store.GetByEnvironmentID(ctx, tc.environmentID)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSpace, result)
			}
		})
	}
}
