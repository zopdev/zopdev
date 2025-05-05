package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gofr.dev/pkg/gofr/container"

	"github.com/DATA-DOG/go-sqlmock"

	"gofr.dev/pkg/gofr"
)

func TestInsertCloudAccount(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	cloudAccount := &CloudAccount{
		Name:            "Test Account",
		Provider:        "GCP",
		ProviderID:      "gcp-project-id",
		ProviderDetails: `{"region":"us-central1"}`,
		Credentials:     map[string]string{"key": "value"},
	}

	testCases := []struct {
		name          string
		cloudAccount  *CloudAccount
		expectedError bool
		mockBehavior  func()
	}{
		{
			name:          "success",
			cloudAccount:  cloudAccount,
			expectedError: false,
			mockBehavior: func() {
				jsonCredentials, _ := json.Marshal(cloudAccount.Credentials)
				mock.SQL.ExpectExec("INSERT INTO cloud_account (name, provider,provider_id,provider_details ,credentials) values(? , ?, ?, ? ,?);").
					WithArgs(cloudAccount.Name, cloudAccount.Provider, cloudAccount.ProviderID, cloudAccount.ProviderDetails, jsonCredentials).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "failure on marshaling credentials",
			cloudAccount: &CloudAccount{Name: "Invalid Account", Provider: "GCP",
				ProviderID: "invalid-id", Credentials: make(chan int)}, // invalid type
			expectedError: true,
			mockBehavior:  func() {},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			_, err := store.InsertCloudAccount(ctx, tc.cloudAccount)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetALLCloudAccounts(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		expectedError bool
		expectedCount int
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockRows := sqlmock.NewRows([]string{"id", "name", "provider", "provider_id", "provider_details", "created_at", "updated_at"}).
					AddRow(1, "Test Account", "GCP", "gcp-project-id", `{"region":"us-central1"}`, time.Now(), time.Now())
				mock.SQL.ExpectQuery(GETALLQUERY).
					WillReturnRows(mockRows)
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "failure on query execution",
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETALLQUERY).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			cloudAccounts, err := store.GetALLCloudAccounts(ctx)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, cloudAccounts)
			} else {
				require.NoError(t, err)
				require.Len(t, cloudAccounts, tc.expectedCount)
			}
		})
	}
}

func TestGetCloudAccountByProvider(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	testCases := []struct {
		name          string
		providerType  string
		providerID    string
		mockBehavior  func()
		expectedError bool
		expectedNil   bool
		expectedName  string
	}{
		{
			name:         "success",
			providerType: "GCP",
			providerID:   "gcp-project-id",
			mockBehavior: func() {
				mockRow := sqlmock.NewRows([]string{"id", "name", "provider", "provider_id", "provider_details", "created_at", "updated_at"}).
					AddRow(1, "Test Account", "GCP", "gcp-project-id", `{"region":"us-central1"}`, time.Now(), time.Now())
				mock.SQL.ExpectQuery(GETBYPROVIDERQUERY).
					WithArgs("GCP", "gcp-project-id").
					WillReturnRows(mockRow)
			},
			expectedError: false,
			expectedNil:   false,
			expectedName:  "Test Account",
		},
		{
			name:         "no rows found",
			providerType: "AWS",
			providerID:   "non-existent-id",
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYPROVIDERQUERY).
					WithArgs("GCP", "non-existent-id").
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedError: true,
			expectedNil:   true,
			expectedName:  "",
		},
		{
			name:         "failure on query execution",
			providerType: "GCP",
			providerID:   "gcp-project-id",
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYPROVIDERQUERY).
					WithArgs("GCP", "gcp-project-id").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			expectedNil:   true,
			expectedName:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			cloudAccount, err := store.GetCloudAccountByProvider(ctx, tc.providerType, tc.providerID)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, cloudAccount)
			} else {
				require.NoError(t, err)

				if tc.expectedNil {
					require.Nil(t, cloudAccount)
				} else {
					require.NotNil(t, cloudAccount)
					require.Equal(t, tc.expectedName, cloudAccount.Name)
				}
			}
		})
	}
}
