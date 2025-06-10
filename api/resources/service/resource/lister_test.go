package resource

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

var errMock = errors.New("mock error")

func TestService_getAllSQLInstances_UnsupportedCloud(t *testing.T) {
	ctx := &gofr.Context{}
	req := CloudDetails{CloudType: "Unknown"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mAWS := NewMockCloudResourceProvider(ctrl)
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

	mockGCP := NewMockCloudResourceProvider(ctrl)
	mockResp := []models.Resource{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
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
				mockGCP.EXPECT().ListResources(ctx, req.Creds, gomock.Any()).Return(mockResp, nil)
			},
		},
		{
			name:   "Error listing resources",
			req:    req,
			expErr: errMock,
			mockCalls: func() {
				mockGCP.EXPECT().ListResources(ctx, req.Creds, gomock.Any()).Return(nil, errMock)
			},
		},
	}

	mAWS := NewMockCloudResourceProvider(ctrl)
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
