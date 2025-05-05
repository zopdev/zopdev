package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestInsertApplication(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	application := &Application{Name: "Test Application"}

	testCases := []struct {
		name          string
		application   *Application
		expectedError bool
		mockBehavior  func()
	}{
		{
			name:          "success",
			application:   application,
			expectedError: false,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(application.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:          "failure on query execution",
			application:   application,
			expectedError: true,
			mockBehavior: func() {
				mock.SQL.ExpectExec(INSERTQUERY).
					WithArgs(application.Name).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			result, err := store.InsertApplication(ctx, tc.application)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, application.Name, result.Name)
			}
		})
	}
}

func TestGetALLApplications(t *testing.T) {
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
				mockRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Test Application", time.Now(), time.Now())
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
			applications, err := store.GetALLApplications(ctx)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, applications)
			} else {
				require.NoError(t, err)
				require.Len(t, applications, tc.expectedCount)
			}
		})
	}
}

func TestGetApplicationByName(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	testCases := []struct {
		name          string
		appName       string
		mockBehavior  func()
		expectedError bool
		expectedNil   bool
		expectedName  string
	}{
		{
			name:    "success",
			appName: "Test Application",
			mockBehavior: func() {
				mockRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Test Application", time.Now(), time.Now())
				mock.SQL.ExpectQuery(GETBYNAMEQUERY).
					WithArgs("Test Application").
					WillReturnRows(mockRow)
			},
			expectedError: false,
			expectedNil:   false,
			expectedName:  "Test Application",
		},
		{
			name:    "no rows found",
			appName: "Non-existent Application",
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYNAMEQUERY).
					WithArgs("Non-existent Application").
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedError: true,
			expectedNil:   true,
			expectedName:  "",
		},
		{
			name:    "failure on query execution",
			appName: "Test Application",
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYNAMEQUERY).
					WithArgs("Test Application").
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
			application, err := store.GetApplicationByName(ctx, tc.appName)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, application)
			} else {
				require.NoError(t, err)

				if tc.expectedNil {
					require.Nil(t, application)
				} else {
					require.NotNil(t, application)
					require.Equal(t, tc.expectedName, application.Name)
				}
			}
		})
	}
}

func TestGetApplicationByID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}

	testCases := []struct {
		name          string
		appID         int
		mockBehavior  func()
		expectedError bool
		expectedNil   bool
		expectedID    int64
	}{
		{
			name:  "success",
			appID: 1,
			mockBehavior: func() {
				mockRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Test Application", time.Now(), time.Now())
				mock.SQL.ExpectQuery(GETBYIDQUERY).
					WithArgs(1).
					WillReturnRows(mockRow)
			},
			expectedError: false,
			expectedNil:   false,
			expectedID:    1,
		},
		{
			name:  "failure on query execution",
			appID: 1,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYIDQUERY).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			expectedNil:   true,
			expectedID:    0,
		},
		{
			name:  "no rows found",
			appID: 1,
			mockBehavior: func() {
				mock.SQL.ExpectQuery(GETBYIDQUERY).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedError: true,
			expectedNil:   true,
			expectedID:    0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			store := New()
			application, err := store.GetApplicationByID(ctx, tc.appID)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, application)
			} else {
				require.NoError(t, err)

				if tc.expectedNil {
					require.Nil(t, application)
				} else {
					require.NotNil(t, application)
					require.Equal(t, tc.expectedID, application.ID)
				}
			}
		})
	}
}
