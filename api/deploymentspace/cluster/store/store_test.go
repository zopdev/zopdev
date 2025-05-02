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

func TestInsertCluster(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	cluster := &Cluster{
		DeploymentSpaceID: 1,
		Identifier:        "test-id",
		Name:              "Test Cluster",
		Region:            "us-central1",
		ProviderID:        "123",
		Provider:          "gcp",
		Namespace:         Namespace{Name: "test-namespace"},
	}

	testCases := []struct {
		name          string
		cluster       *Cluster
		expectedError bool
		mockBehavior  func()
	}{
		{
			name:          "success",
			cluster:       cluster,
			expectedError: false,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(cluster.DeploymentSpaceID, cluster.Identifier, cluster.Name, cluster.Region,
						cluster.ProviderID, cluster.Provider, cluster.Namespace.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:          "error inserting cluster",
			cluster:       cluster,
			expectedError: true,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(cluster.DeploymentSpaceID, cluster.Identifier, cluster.Name, cluster.Region,
						cluster.ProviderID, cluster.Provider, cluster.Namespace.Name).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			result, err := store.Insert(ctx, tc.cluster)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, int64(1), result.ID)
			}
		})
	}
}

func TestGetClusterByDeploymentSpaceID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	expectedCluster := &Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Identifier:        "test-id",
		Name:              "Test Cluster",
		Region:            "us-central1",
		ProviderID:        "123",
		Provider:          "gcp",
		Namespace:         Namespace{Name: "test-namespace"},
		CreatedAt:         "2023-12-11T00:00:00Z",
		UpdatedAt:         "2023-12-11T00:00:00Z",
	}

	testCases := []struct {
		name              string
		deploymentSpaceID int
		mockBehavior      func()
		expectedError     bool
		expectedCluster   *Cluster
	}{
		{
			name:              "success",
			deploymentSpaceID: 1,
			expectedError:     false,
			mockBehavior: func() {
				mockRow := sqlmock.NewRows([]string{"id", "deployment_space_id", "identifier", "name", "region", "provider_id",
					"provider", "namespace", "created_at", "updated_at"}).
					AddRow(1, 1, "test-id", "Test Cluster", "us-central1", 123, "gcp", "test-namespace",
						"2023-12-11T00:00:00Z", "2023-12-11T00:00:00Z")
				mock.SQL.ExpectQuery(GETQUERY).
					WithArgs(1).
					WillReturnRows(mockRow)
			},
			expectedCluster: expectedCluster,
		},
		{
			name:              "no cluster found",
			deploymentSpaceID: 1,
			expectedError:     true,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETQUERY).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(nil)) // No rows returned
			},
			expectedCluster: nil,
		},
		{
			name:              "error on query execution",
			deploymentSpaceID: 1,
			expectedError:     true,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETQUERY).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCluster: nil,
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			result, err := store.GetByDeploymentSpaceID(ctx, tc.deploymentSpaceID)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedCluster, result)
			}
		})
	}
}
