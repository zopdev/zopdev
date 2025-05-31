package resourcegroup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zopdev/zopdev/api/resources/models"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetAllResourceGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mClient := NewMockHTTPClient(ctrl)
	mGCP := NewMockGCP(ctrl)

	s := New(mClient, mGCP)

	ctx := gofr.NewContext(nil, nil, nil)

	ca := &models.CloudAccount{
		ID:   123,
		Type: "gcp",
	}

	req := &models.CloudCredentials{
		Creds: "test-creds",
	}

	mClient.EXPECT().GetCloudCredentials(ctx, int64(123)).Return(ca, nil)
	mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
		Return(nil, nil)

	// Mocking the response for GetAllInstances
	mockResp := []models.Resource{
		{
			ID:           1,
			Name:         "test-instance",
			Type:         "sql",
			CloudAccount: *ca,
			Region:       "us-central1",
			CreationTime: "2023-01-01T00:00:00Z",
			Status:       "running",
			UID:          "uid-123",
			Settings:     models.Settings{"setting1": "value1"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mGCP.EXPECT().GetAllInstances(ctx, ca.ID).Return(mockResp, nil)

	res, err := s.GetAllResourceGroups(ctx, int64(123), []string{"sql"})

	assert.NoError(t, err)
	assert.Equal(t, mockResp, res)
}
