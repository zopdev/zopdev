package service

import (
	"context"
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
