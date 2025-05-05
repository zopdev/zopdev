package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/cloudaccounts/service"
	"github.com/zopdev/zopdev/api/deploymentspace"
	clusterStore "github.com/zopdev/zopdev/api/deploymentspace/cluster/store"
	"github.com/zopdev/zopdev/api/deploymentspace/store"
	"github.com/zopdev/zopdev/api/provider"
)

var (
	errTest        = errors.New("service error")
	errStore       = errors.New("store error")
	errClusterSvc  = errors.New("cluster service error")
	errCloudAccSvc = errors.New("cloud account service error")
	errProviderSvc = errors.New("provider service error")
)

//nolint:funlen // test function
func TestService_AddDeploymentSpace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockDeploymentSpaceStore(ctrl)
	mockClusterService := deploymentspace.NewMockDeploymentEntity(ctrl)

	ctx := &gofr.Context{}

	deploymentSpace := &DeploymentSpace{
		CloudAccount: CloudAccount{
			ID:         1,
			Provider:   "aws",
			ProviderID: "provider-123",
		},
		Type: Type{Name: "test-type"},
		DeploymentSpace: map[string]interface{}{
			"key": "value",
		},
	}

	mockCluster := clusterStore.Cluster{
		DeploymentSpaceID: 1,
		Provider:          "aws",
		ProviderID:        "provider-123",
	}

	mockDeploymentSpace := &store.DeploymentSpace{ID: 1}

	testCases := []struct {
		name          string
		mockBehavior  func()
		input         *DeploymentSpace
		envID         int
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, gomock.Any()).
					Return(nil, nil)
				mockClusterService.EXPECT().
					DuplicateCheck(ctx, gomock.Any()).
					Return(nil, nil) // No duplicate found
				mockStore.EXPECT().
					Insert(ctx, gomock.Any()).
					Return(mockDeploymentSpace, nil)
				mockClusterService.EXPECT().
					Add(ctx, gomock.Any()).
					Return(&mockCluster, nil)
			},
			input:         deploymentSpace,
			envID:         1,
			expectedError: nil,
		},
		{
			name: "duplicate cluster error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, gomock.Any()).
					Return(nil, nil)
				mockClusterService.EXPECT().
					DuplicateCheck(ctx, gomock.Any()).
					Return(nil, http.ErrorEntityAlreadyExist{}) // Duplicate found
			},
			input:         deploymentSpace,
			envID:         1,
			expectedError: http.ErrorEntityAlreadyExist{},
		},
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, gomock.Any()).
					Return(nil, nil)
				mockClusterService.EXPECT().
					DuplicateCheck(ctx, gomock.Any()).
					Return(nil, nil) // No duplicate found
				mockStore.EXPECT().
					Insert(ctx, gomock.Any()).
					Return(nil, errTest)
			},
			input:         deploymentSpace,
			envID:         1,
			expectedError: errTest,
		},
		{
			name: "cluster service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, gomock.Any()).
					Return(nil, nil)
				mockClusterService.EXPECT().
					DuplicateCheck(ctx, gomock.Any()).
					Return(nil, nil) // No duplicate found
				mockStore.EXPECT().
					Insert(ctx, gomock.Any()).
					Return(mockDeploymentSpace, nil)
				mockClusterService.EXPECT().
					Add(ctx, gomock.Any()).
					Return(nil, errTest)
			},
			input:         deploymentSpace,
			envID:         1,
			expectedError: errTest,
		},
		{
			name:         "invalid request body",
			mockBehavior: func() {},
			input: &DeploymentSpace{
				CloudAccount:    CloudAccount{},
				Type:            Type{},
				DeploymentSpace: nil, // Invalid DeploymentEntity
			},
			envID:         1,
			expectedError: http.ErrorInvalidParam{Params: []string{"body"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := New(mockStore, mockClusterService, nil, nil)
			_, err := svc.Add(ctx, tc.input, tc.envID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_FetchDeploymentSpace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockDeploymentSpaceStore(ctrl)
	mockClusterService := deploymentspace.NewMockDeploymentEntity(ctrl)

	ctx := &gofr.Context{}

	mockDeploymentSpace := &store.DeploymentSpace{
		ID:             1,
		CloudAccountID: 1,
		EnvironmentID:  1,
		Type:           "test-type",
	}

	mockCluster := clusterStore.Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Name:              "test-cluster",
	}

	testCases := []struct {
		name          string
		mockBehavior  func()
		envID         int
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(mockDeploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(mockCluster, nil)
			},
			envID:         1,
			expectedError: nil,
		},
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(nil, errTest)
			},
			envID:         1,
			expectedError: errTest,
		},
		{
			name: "no cluster found",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(mockDeploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(nil, sql.ErrNoRows)
			},
			envID:         1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := New(mockStore, mockClusterService, nil, nil)
			resp, err := svc.Fetch(ctx, tc.envID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestGetDeploymentSpaceArgs(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         *gofr.Context
		cluster     *store.Cluster
		credentials interface{}
		expectedCl  *provider.Cluster
		expectedCa  *provider.CloudAccount
		expectedNs  string
	}{
		{
			name: "valid input",
			ctx:  &gofr.Context{},
			cluster: &store.Cluster{
				Name:       "test-cluster",
				Region:     "us-west-1",
				Provider:   "aws",
				ProviderID: "provider-123",
				Namespace:  store.Namespace{Name: "test-namespace"},
			},
			credentials: "test-credentials",
			expectedCl: &provider.Cluster{
				Name:   "test-cluster",
				Region: "us-west-1",
			},
			expectedCa: &provider.CloudAccount{
				Provider:   "aws",
				ProviderID: "provider-123",
			},
			expectedNs: "test-namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cl, _ := getClusterCloudAccount(tc.cluster)

			require.Equal(t, tc.expectedCl, cl)
		})
	}
}

//nolint:funlen //test function
func TestService_GetServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockDeploymentSpaceStore(ctrl)
	mockClusterService := deploymentspace.NewMockDeploymentEntity(ctrl)
	mockCloudAccountService := service.NewMockCloudAccountService(ctrl)
	mockProviderService := provider.NewMockProvider(ctrl)

	ctx := &gofr.Context{}

	deploymentSpace := &store.DeploymentSpace{
		ID:             1,
		CloudAccountID: 1,
		EnvironmentID:  1,
		Type:           "test-type",
	}

	cluster := store.Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Name:              "test-cluster",
		Region:            "us-west-1",
		Provider:          "aws",
		ProviderID:        "provider-123",
		Namespace:         store.Namespace{Name: "test-namespace"},
	}

	credentials := "test-credentials"

	testCases := []struct {
		name          string
		mockBehavior  func()
		envID         int
		expectedError error
		expectedResp  any
	}{
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(nil, errStore)
			},
			envID:         1,
			expectedError: errStore,
			expectedResp:  nil,
		},
		{
			name: "cluster service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(nil, errClusterSvc)
			},
			envID:         1,
			expectedError: errClusterSvc,
			expectedResp:  nil,
		},
		{
			name: "cloud account service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(cluster, nil)
				mockCloudAccountService.EXPECT().
					FetchCredentials(ctx, int64(1)).
					Return(nil, errCloudAccSvc)
			},
			envID:         1,
			expectedError: errCloudAccSvc,
			expectedResp:  nil,
		},
		{
			name: "provider service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(cluster, nil)
				mockCloudAccountService.EXPECT().
					FetchCredentials(ctx, int64(1)).
					Return(credentials, nil)
				mockProviderService.EXPECT().
					ListServices(ctx, &provider.Cluster{
						Name:   "test-cluster",
						Region: "us-west-1",
					}, &provider.CloudAccount{
						Provider:   "aws",
						ProviderID: "provider-123",
					}, credentials, "test-namespace").
					Return(nil, errProviderSvc)
			},
			envID:         1,
			expectedError: errProviderSvc,
			expectedResp:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := New(mockStore, mockClusterService, mockCloudAccountService, mockProviderService)
			resp, err := svc.GetServices(ctx, tc.envID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp, resp)
			}
		})
	}
}

//nolint:funlen //test function
func TestService_GetDeployments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockDeploymentSpaceStore(ctrl)
	mockClusterService := deploymentspace.NewMockDeploymentEntity(ctrl)
	mockCloudAccountService := service.NewMockCloudAccountService(ctrl)
	mockProviderService := provider.NewMockProvider(ctrl)

	ctx := &gofr.Context{}

	deploymentSpace := &store.DeploymentSpace{
		ID:             1,
		CloudAccountID: 1,
		EnvironmentID:  1,
		Type:           "test-type",
	}

	cluster := store.Cluster{
		ID:                1,
		DeploymentSpaceID: 1,
		Name:              "test-cluster",
		Region:            "us-west-1",
		Provider:          "aws",
		ProviderID:        "provider-123",
		Namespace:         store.Namespace{Name: "test-namespace"},
	}

	credentials := "test-credentials"

	testCases := []struct {
		name          string
		mockBehavior  func()
		envID         int
		expectedError error
		expectedResp  any
	}{
		{
			name: "store layer error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(nil, errStore)
			},
			envID:         1,
			expectedError: errStore,
			expectedResp:  nil,
		},
		{
			name: "cluster service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(nil, errClusterSvc)
			},
			envID:         1,
			expectedError: errClusterSvc,
			expectedResp:  nil,
		},
		{
			name: "cloud account service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(cluster, nil)
				mockCloudAccountService.EXPECT().
					FetchCredentials(ctx, int64(1)).
					Return(nil, errCloudAccSvc)
			},
			envID:         1,
			expectedError: errCloudAccSvc,
			expectedResp:  nil,
		},
		{
			name: "provider service error",
			mockBehavior: func() {
				mockStore.EXPECT().
					GetByEnvironmentID(ctx, 1).
					Return(deploymentSpace, nil)
				mockClusterService.EXPECT().
					FetchByDeploymentSpaceID(ctx, 1).
					Return(cluster, nil)
				mockCloudAccountService.EXPECT().
					FetchCredentials(ctx, int64(1)).
					Return(credentials, nil)
				mockProviderService.EXPECT().
					ListDeployments(ctx, &provider.Cluster{
						Name:   "test-cluster",
						Region: "us-west-1",
					}, &provider.CloudAccount{
						Provider:   "aws",
						ProviderID: "provider-123",
					}, credentials, "test-namespace").
					Return(nil, errProviderSvc)
			},
			envID:         1,
			expectedError: errProviderSvc,
			expectedResp:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			svc := New(mockStore, mockClusterService, mockCloudAccountService, mockProviderService)
			resp, err := svc.GetDeployments(ctx, tc.envID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp, resp)
			}
		})
	}
}
