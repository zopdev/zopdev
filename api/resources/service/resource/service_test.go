package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

func TestService_SyncResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mClient := NewMockHTTPClient(ctrl)
	mStore := NewMockStore(ctrl)
	mGCP := NewMockGCPClient(ctrl)
	mAWS := NewMockAWSClient(ctrl)

	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{Context: context.Background(), Container: mockContainer}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	mockCreds := &google.Credentials{ProjectID: "test-project"}

	s := New(mGCP, mAWS, mClient, mStore)

	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}
	mockInst := []models.Resource{
		{Name: "sql-instance-1", UID: "zopdev/sql-instance-1", Type: "SQL", Status: "RUNNING"},
		{Name: "sql-instance-2", UID: "zopdev/sql-instance-2", Type: "SQL", Status: "SUSPENDED"},
	}
	mStrResp := []models.Resource{
		{ID: 1, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
			Name: "sql-instance-1", Type: string(SQL),
			UID: "zopdev/sql-instance-1", Status: "RUNNING"},
		{ID: 3, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
			Name: "sql-instance-2", Type: string(SQL),
			UID: "zopdev/sql-instance-3", Status: "SUSPENDED"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockInst,
	}

	testCases := []struct {
		name      string
		id        int64
		resources []string
		expErr    error
		expResp   []models.Resource
		mockCalls func()
	}{
		{
			name:      "Sync all resources",
			id:        123,
			resources: []string{},
			expErr:    nil,
			expResp:   mStrResp,
			mockCalls: func() {
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
				gomock.InOrder(
					mStore.EXPECT().GetResources(gomock.Any(), int64(123), nil).
						Return([]models.Resource{
							{ID: 1, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
								Name: "sql-instance-1", Type: string(SQL), UID: "zopdev/sql-instance-1"},
							{ID: 2, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
								Name: "sql-instance-3", Type: string(SQL), UID: "zopdev/sql-instance-3"},
						}, nil),
					mStore.EXPECT().UpdateStatus(gomock.Any(), RUNNING, int64(1)).Return(nil),
					mStore.EXPECT().InsertResource(gomock.Any(), gomock.Any()).Return(nil),
					mStore.EXPECT().RemoveResource(gomock.Any(), gomock.Any()).Return(nil),
					mStore.EXPECT().GetResources(gomock.Any(), int64(123), nil).Return([]models.Resource{
						{ID: 1, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
							Name: "sql-instance-1", Type: string(SQL), UID: "zopdev/sql-instance-1", Status: "RUNNING"},
						{ID: 3, CloudAccount: models.CloudAccount{ID: 123, Type: string(GCP)},
							Name: "sql-instance-2", Type: string(SQL), UID: "zopdev/sql-instance-3", Status: "SUSPENDED"},
					}, nil),
					mStore.EXPECT().GetResources(gomock.Any(), int64(123), nil).Return(mStrResp, nil).AnyTimes(),
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.SyncResources(ctx, tc.id)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_SyncResources_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mClient := NewMockHTTPClient(ctrl)
	mGCP := NewMockGCPClient(ctrl)
	ct, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{Context: context.Background(), Container: ct}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	s := New(mGCP, nil, mClient, nil)
	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}

	testCases := []struct {
		name      string
		id        int64
		resources []string
		expErr    error
		expResp   []models.Resource
		mockCalls func()
	}{
		{
			name:      "error getting SQL resources",
			id:        123,
			resources: []string{string(SQL)},
			expErr:    errMock,
			mockCalls: func() {
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "error from GetCloudCredentials",
			id:     123,
			expErr: errMock,
			mockCalls: func() {
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(nil, errMock)
			},
		},
		{
			name:      "error from GetAllInstances",
			id:        123,
			resources: []string{},
			expErr:    errMock,
			mockCalls: func() {
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.SyncResources(ctx, tc.id)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_ChangeState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mClient := NewMockHTTPClient(ctrl)
	mStore := NewMockStore(ctrl)
	mGCP := NewMockGCPClient(ctrl)
	mAWS := NewMockAWSClient(ctrl)
	ctx := &gofr.Context{Context: context.Background()}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockStopper := &mockSQLClient{}
	s := New(mGCP, mAWS, mClient, mStore)

	testCases := []struct {
		name      string
		input     ResourceDetails
		expErr    error
		mockCalls func()
	}{
		{
			name:  "Success - Start",
			input: ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(&models.Resource{ID: 1, CloudAccount: models.CloudAccount{ID: 1}, Status: STOPPED}, nil)
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).
					Return(ca, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
				mStore.EXPECT().UpdateStatus(ctx, RUNNING, int64(1)).
					Return(nil)
			},
		},
		{
			name:  "Success - Suspend",
			input: ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(&models.Resource{ID: 1, CloudAccount: models.CloudAccount{ID: 1}, Status: RUNNING}, nil)
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
				mStore.EXPECT().UpdateStatus(ctx, STOPPED, int64(1)).
					Return(nil)
			},
		},
		{
			name:  "Success - Resource already in desired state",
			input: ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(&models.Resource{ID: 1, CloudAccount: models.CloudAccount{ID: 1}, Status: RUNNING}, nil)
			},
		},
		{
			name:   "Error - GetResourceByID",
			input:  ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			expErr: errMock,
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error - GetCloudCredentials",
			input:  ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			expErr: errMock,
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(&models.Resource{ID: 1, CloudAccount: models.CloudAccount{ID: 1}, Status: STOPPED}, nil)
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error - Invalid Type",
			input:  ResourceDetails{ID: 1, CloudAccID: 123, Name: "test-instance", Type: "invalid"},
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"req.Type"}},
			mockCalls: func() {
				mStore.EXPECT().GetResourceByID(ctx, int64(1)).
					Return(&models.Resource{ID: 1, CloudAccount: models.CloudAccount{ID: 1}, Status: RUNNING}, nil)
				mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := s.ChangeState(ctx, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}
