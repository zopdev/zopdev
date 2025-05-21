package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

func InitializeTests(t *testing.T) (*gofr.Context, *gomock.Controller, *container.Mocks, *MockGCPClient) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockContainer, mock := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))

	ctx := &gofr.Context{
		Container: mockContainer,
	}
	mGCP := NewMockGCPClient(ctrl)

	return ctx, ctrl, mock, mGCP
}

func TestService_GetResources(t *testing.T) {
	ctx, ctrl, mock, mGCP := InitializeTests(t)
	defer ctrl.Finish()

	s := New(mGCP)
	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockResp := []models.Instance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}

	testCases := []struct {
		name      string
		id        int64
		resources []string
		expErr    error
		expResp   []models.Instance
		mockCalls func()
	}{
		{
			name:      "Get all resources",
			id:        123,
			resources: []string{},
			expResp:   mockResp,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
		{
			name:      "Get SQL resources",
			id:        123,
			resources: []string{string(SQL)},
			expResp:   mockResp,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.GetResources(ctx, tc.id, tc.resources)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_GetResources_Errors(t *testing.T) {
	ctx, ctrl, mock, mGCP := InitializeTests(t)
	defer ctrl.Finish()

	s := New(mGCP)
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
		expResp   []models.Instance
		mockCalls func()
	}{
		{
			name:      "error getting SQL resources",
			id:        123,
			resources: []string{string(SQL)},
			expErr:    errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "error from GetCloudCredentials",
			id:     123,
			expErr: errMock,
			mockCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:      "error from GetAllInstances",
			id:        123,
			resources: []string{},
			expErr:    errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)

				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:      "unsupported resource",
			id:        123,
			resources: []string{"unsupported"},
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}},
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.GetResources(ctx, tc.id, tc.resources)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_getAllSQLInstances_UnsupportedCloud(t *testing.T) {
	ctx := &gofr.Context{}
	req := CloudDetails{CloudType: "AWS"}

	s := New(nil)
	instances, err := s.getAllSQLInstances(ctx, req)

	assert.Nil(t, instances)
	require.Error(t, err)
	assert.Equal(t, gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}}, err)
}

func TestService_getAllSQLInstances_GCP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}
	ctx := &gofr.Context{
		Context: context.Background(),
	}

	mockGCP := NewMockGCPClient(ctrl)
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockResp := []models.Instance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}

	testCases := []struct {
		name      string
		req       CloudDetails
		expResp   []models.Instance
		expErr    error
		mockCalls func()
	}{
		{
			name:    "Success",
			req:     req,
			expResp: mockResp,
			expErr:  nil,
			mockCalls: func() {
				mockGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mockGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
		{
			name:   "Error creating credentials",
			req:    req,
			expErr: errMock,
			mockCalls: func() {
				mockGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "Error creating SQL instance lister",
			req:    req,
			expErr: errMock,
			mockCalls: func() {
				mockGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mockGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error getting SQL instances",
			req:    req,
			expErr: errMock,
			mockCalls: func() {
				mockLister.isError = true

				mockGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mockGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
	}

	s := New(mockGCP)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			instances, err := s.getAllSQLInstances(ctx, req)

			assert.Equal(t, tc.expResp, instances)
			assert.Equal(t, tc.expErr, err)
		})
	}
}
