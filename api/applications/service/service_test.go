package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"

	envStore "github.com/zopdev/zopdev/api/environments/store"

	"github.com/zopdev/zopdev/api/applications/store"
	"github.com/zopdev/zopdev/api/environments/service"
	"gofr.dev/pkg/gofr/http"
)

var (
	errTest = errors.New("service error")
)

func TestService_AddApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockApplicationStore(ctrl)
	mockEvironmentService := service.NewMockEnvironmentService(ctrl)
	ctx := &gofr.Context{}

	application := &store.Application{
		Name: "Test Application",
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		input         *store.Application
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByName(ctx, "Test Application").
					Return(nil, sql.ErrNoRows)
				mockStore.EXPECT().
					InsertApplication(ctx, application).
					Return(application, nil)
				mockStore.EXPECT().
					InsertEnvironment(ctx, gomock.Any()).
					Return(&store.Environment{ID: 1}, nil).Times(1)
			},
			input:         application,
			expectedError: nil,
		},
		{
			name: "application already exists",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByName(ctx, "Test Application").
					Return(application, nil)
			},
			input:         application,
			expectedError: http.ErrorEntityAlreadyExist{},
		},
		{
			name: "error fetching application by name",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByName(ctx, "Test Application").
					Return(nil, errTest)
			},
			input:         application,
			expectedError: errTest,
		},
		{
			name: "error inserting application",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByName(ctx, "Test Application").
					Return(nil, sql.ErrNoRows)
				mockStore.EXPECT().
					InsertApplication(ctx, application).
					Return(nil, errTest)
			},
			input:         application,
			expectedError: errTest,
		},
		{
			name: "error inserting environment",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByName(ctx, "Test Application").
					Return(nil, sql.ErrNoRows)
				mockStore.EXPECT().
					InsertApplication(ctx, application).
					Return(application, nil)
				mockStore.EXPECT().
					InsertEnvironment(ctx, gomock.Any()).
					Return(nil, errTest).Times(1)
			},
			input:         application,
			expectedError: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			appService := New(mockStore, mockEvironmentService)
			_, err := appService.AddApplication(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_FetchAllApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockApplicationStore(ctrl)
	mockEvironmentService := service.NewMockEnvironmentService(ctrl)

	ctx := &gofr.Context{}

	expectedApplications := []store.Application{
		{
			ID:        1,
			Name:      "Test Application",
			CreatedAt: "2023-12-11T00:00:00Z",
		},
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetALLApplications(ctx).
					Return(expectedApplications, nil)
				mockEvironmentService.EXPECT().
					FetchAll(ctx, 1).
					Return([]envStore.Environment{{ID: 1, Name: "default", Level: 1}}, nil)
			},
			expectedError: nil,
		},
		{
			name: "error fetching applications",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetALLApplications(ctx).
					Return(nil, errTest)
			},
			expectedError: errTest,
		},
		{
			name: "error fetching environments for application",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetALLApplications(ctx).
					Return(expectedApplications, nil)
				mockEvironmentService.EXPECT().
					FetchAll(ctx, 1).
					Return(nil, errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			appService := New(mockStore, mockEvironmentService)
			applications, err := appService.FetchAllApplications(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedApplications, applications)
			}
		})
	}
}

func TestService_GetApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockApplicationStore(ctrl)
	mockEvironmentService := service.NewMockEnvironmentService(ctrl)

	ctx := &gofr.Context{}

	expectedApplication := store.Application{
		ID:        1,
		Name:      "Test Application",
		CreatedAt: "2023-12-11T00:00:00Z",
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		expectedError error
		expectedApp   *store.Application
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByID(ctx, 1).
					Return(&expectedApplication, nil)
				mockEvironmentService.EXPECT().
					FetchAll(ctx, 1).
					Return([]envStore.Environment{{ID: 1, Name: "default", Level: 1}}, nil)
			},
			expectedError: nil,
			expectedApp:   &expectedApplication,
		},
		{
			name: "error fetching application by ID",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByID(ctx, 1).
					Return(nil, errTest)
			},
			expectedError: errTest,
			expectedApp:   nil,
		},
		{
			name: "error fetching environments for application",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetApplicationByID(ctx, 1).
					Return(&expectedApplication, nil)
				mockEvironmentService.EXPECT().
					FetchAll(ctx, 1).
					Return(nil, errTest)
			},
			expectedError: errTest,
			expectedApp:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			appService := New(mockStore, mockEvironmentService)
			application, err := appService.GetApplication(ctx, 1)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedApp, application)
			}
		})
	}
}
