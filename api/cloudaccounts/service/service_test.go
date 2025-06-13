package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/cloudaccounts/store"
	"github.com/zopdev/zopdev/api/provider"
)

var (
	errTest = errors.New("service error")
)

// TODO: add tests for AWS as well once the structure is fixed as the tests would need to be rewritten.
func TestService_AddGCPCloudAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockCloudAccountStore(ctrl)
	mockProvider := provider.NewMockProvider(ctrl)

	ctx := &gofr.Context{}

	gcpCloudAccount := &store.CloudAccount{
		Name:        "Test Account",
		Provider:    "GCP",
		ProviderID:  "test-project-id",
		Credentials: map[string]string{"project_id": "test-project-id"},
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		input         *store.CloudAccount
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetCloudAccountByProvider(ctx, "GCP", "test-project-id").
					Return(nil, nil)
				mockStore.EXPECT().
					InsertCloudAccount(ctx, gcpCloudAccount).
					Return(gcpCloudAccount, nil)
			},
			input:         gcpCloudAccount,
			expectedError: nil,
		},
		{
			name: "duplicate account",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetCloudAccountByProvider(ctx, "GCP", "test-project-id").
					Return(gcpCloudAccount, nil)
			},
			input:         gcpCloudAccount,
			expectedError: http.ErrorEntityAlreadyExist{},
		},
		{
			name: "error getting account",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetCloudAccountByProvider(ctx, "GCP", "test-project-id").
					Return(nil, errTest)
			},
			input:         gcpCloudAccount,
			expectedError: errTest,
		},
		{
			name:         "invalid credentials for GCP",
			mockBehavior: func() {},
			input: &store.CloudAccount{
				Name:        "Invalid Account",
				Provider:    "GCP",
				Credentials: map[string]string{},
			},
			expectedError: http.ErrorInvalidParam{Params: []string{"credentials"}},
		},
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetCloudAccountByProvider(ctx, "GCP", "test-project-id").
					Return(nil, nil)
				mockStore.EXPECT().
					InsertCloudAccount(ctx, gcpCloudAccount).
					Return(nil, errTest)
			},
			input:         gcpCloudAccount,
			expectedError: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			service := New(mockStore, mockProvider, "")
			_, err := service.AddCloudAccount(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_FetchAllCloudAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockCloudAccountStore(ctrl)
	mockProvider := provider.NewMockProvider(ctrl)

	ctx := &gofr.Context{}

	expectedAccounts := []store.CloudAccount{
		{
			ID:              1,
			Name:            "Test Account",
			Provider:        "GCP",
			ProviderID:      "test-project-id",
			ProviderDetails: `{"region":"us-central1"}`,
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
					GetALLCloudAccounts(ctx).
					Return(expectedAccounts, nil)
			},
			expectedError: nil,
		},
		{
			name: "error fetching accounts",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetALLCloudAccounts(ctx).
					Return(nil, errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			service := New(mockStore, mockProvider, "")
			_, err := service.FetchAllCloudAccounts(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
