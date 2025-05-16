package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zopdev/zopdev/api/resources/providers/models"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"testing"

	"go.uber.org/mock/gomock"
)

var errMock = errors.New("mock error")

type mockInstanceLister struct {
	isError   bool
	instances []models.SQLInstance
}

func (m *mockInstanceLister) GetAllInstances(_ *gofr.Context, _ string) ([]models.SQLInstance, error) {
	if m.isError {
		return nil, errMock
	}

	return m.instances, nil
}

func TestService_GetAllSQLInstances_UnsupportedCloud(t *testing.T) {
	ctx := &gofr.Context{}
	req := Request{CloudType: UNSUPPORTED}

	s := New(nil)
	instances, err := s.GetAllSQLInstances(ctx, req)

	assert.Nil(t, instances)
	assert.Error(t, err)
	assert.Equal(t, gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}}, err)
}

func TestService_GetAllSQLInstances_GCP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := Request{
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
	mockResp := []models.SQLInstance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	mockLister := &mockInstanceLister{
		isError:   false,
		instances: mockResp,
	}

	testCases := []struct {
		name      string
		req       Request
		expResp   []models.SQLInstance
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
				mockGCP.EXPECT().NewSQLInstanceLister(ctx, option.WithCredentials(mockCreds)).
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
				mockGCP.EXPECT().NewSQLInstanceLister(ctx, option.WithCredentials(mockCreds)).
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
				mockGCP.EXPECT().NewSQLInstanceLister(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
	}

	s := New(mockGCP)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			instances, err := s.GetAllSQLInstances(ctx, req)

			assert.Equal(t, tc.expResp, instances)
			assert.Equal(t, tc.expErr, err)
		})
	}
}
