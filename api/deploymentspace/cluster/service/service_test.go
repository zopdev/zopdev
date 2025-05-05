package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zopdev/zopdev/api/deploymentspace/cluster/service"
	"github.com/zopdev/zopdev/api/deploymentspace/cluster/store"
	"gofr.dev/pkg/gofr"
)

var (
	errNamespaceAlreadyInUSe = errors.New("namespace already in use")
	errTest                  = errors.New("service error")
)

func TestService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockClusterStore(ctrl)
	ctx := &gofr.Context{}

	// Mock input data
	cluster := &store.Cluster{
		DeploymentSpaceID: 1,
		Provider:          "aws",
		ProviderID:        "provider-123",
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		input         interface{}
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					Insert(ctx, gomock.Any()).
					Return(cluster, nil)
			},
			input:         cluster,
			expectedError: nil,
		},
		{
			name: "error in Insert",
			mockBehavior: func() {
				mockStore.EXPECT().
					Insert(ctx, gomock.Any()).
					Return(nil, errTest)
			},
			input:         cluster,
			expectedError: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := service.New(mockStore)
			_, err := svc.Add(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_FetchByDeploymentSpaceID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockClusterStore(ctrl)
	ctx := &gofr.Context{}

	expectedCluster := &store.Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Identifier:        "cluster-1",
		Provider:          "aws",
		ProviderID:        "provider-123",
	}

	testCases := []struct {
		name            string
		mockBehavior    func()
		inputID         int
		expectedError   error
		expectedCluster *store.Cluster
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByDeploymentSpaceID(ctx, 1).
					Return(expectedCluster, nil)
			},
			inputID:         1,
			expectedError:   nil,
			expectedCluster: expectedCluster,
		},
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByDeploymentSpaceID(ctx, 1).
					Return(nil, errTest)
			},
			inputID:         1,
			expectedError:   errTest,
			expectedCluster: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := service.New(mockStore)
			result, err := svc.FetchByDeploymentSpaceID(ctx, tc.inputID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedCluster, result)
			}
		})
	}
}

func TestService_DuplicateCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockClusterStore(ctrl)
	ctx := &gofr.Context{}

	mockCluster := &store.Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Provider:          "aws",
		ProviderID:        "provider-123",
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		input         any
		expectedError error
		expectedResp  interface{}
	}{
		{
			name: "success - no duplicate found",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByCluster(ctx, mockCluster).
					Return(nil, sql.ErrNoRows)
			},
			input:         mockCluster,
			expectedError: nil,
			expectedResp:  nil,
		},
		{
			name: "error during GetByCluster",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByCluster(ctx, mockCluster).
					Return(nil, errTest)
			},
			input:         mockCluster,
			expectedError: errTest,
			expectedResp:  nil,
		},
		{
			name: "duplicate cluster found",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByCluster(ctx, mockCluster).
					Return(mockCluster, nil)
			},
			input:         mockCluster,
			expectedError: errNamespaceAlreadyInUSe,
			expectedResp:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := service.New(mockStore)
			resp, err := svc.DuplicateCheck(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedResp, resp)
		})
	}
}
