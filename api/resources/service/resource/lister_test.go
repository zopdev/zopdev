package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestService_getAllSQLInstances_UnsupportedCloud(t *testing.T) {
	ctx := &gofr.Context{}
	req := CloudDetails{CloudType: "Unknown"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mAWS := NewMockAWSClient(ctrl)
	s := New(nil, mAWS, nil, nil)
	instances, err := s.getAllSQLInstances(ctx, req)

	assert.Nil(t, instances)
	require.NoError(t, err)
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
	mockResp := []models.Resource{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}

	testCases := []struct {
		name      string
		req       CloudDetails
		expResp   []models.Resource
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

	mAWS := NewMockAWSClient(ctrl)
	s := New(mockGCP, mAWS, nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			instances, err := s.getAllSQLInstances(ctx, req)

			assert.Equal(t, tc.expResp, instances)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestService_bSearch(t *testing.T) {
	res := []models.Resource{
		{ID: 1, UID: "zopdev-test/mysql01"},
		{ID: 2, UID: "zopdev-test/mysql02"},
		{ID: 3, UID: "zopdev-test/pgs1l01"},
	}

	idx, found := bSearch(res, "zopdev-test/mysql02")

	assert.True(t, found)
	assert.Equal(t, 1, idx)

	idx, found = bSearch(res, "zopdev-test/mysql03")
	assert.False(t, found)
	assert.Equal(t, -1, idx)
}
