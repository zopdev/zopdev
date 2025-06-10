package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func TestService_SyncCron(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cnt, mocks := container.NewMockContainer(t)

	mGCP := NewMockCloudResourceProvider(ctrl)
	mHTTP := NewMockHTTPClient(ctrl)
	mStore := NewMockStore(ctrl)
	mAWS := NewMockCloudResourceProvider(ctrl)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: cnt,
	}
	mockResp := []models.Resource{
		{ID: 1, Name: "sql-instance-1", UID: "zop/sql1", Status: RUNNING},
		{ID: 2, Name: "sql-instance-2", UID: "zop/sql2", Status: STOPPED},
	}

	service := New(mGCP, mAWS, mHTTP, mStore)

	// Test successful sync
	t.Run("successful sync", func(t *testing.T) {
		// mock expectations
		mHTTP.EXPECT().GetAllCloudAccounts(ctx).
			Return([]client.CloudAccount{{ID: 1, Provider: "GCP"}, {ID: 2, Provider: "Unknown"}}, nil)
		mHTTP.EXPECT().GetCloudCredentials(ctx, int64(1)).
			Return(&client.CloudAccount{ID: 1, Provider: "GCP"}, nil)
		mHTTP.EXPECT().GetCloudCredentials(ctx, int64(2)).
			Return(&client.CloudAccount{ID: 2, Provider: "Unknown"}, nil)

		mGCP.EXPECT().ListResources(ctx, gomock.Any(), gomock.Any()).Return(mockResp, nil).AnyTimes()
		mAWS.EXPECT().ListResources(ctx, gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		mStore.EXPECT().GetResources(ctx, int64(1), nil).
			Return(mockResp, nil).AnyTimes()
		mStore.EXPECT().GetResources(ctx, int64(2), nil).
			Return(nil, nil).AnyTimes()

		// Update expectations to match actual behavior
		mStore.EXPECT().UpdateStatus(ctx, gomock.Any(), int64(1)).Return(nil).AnyTimes()
		mStore.EXPECT().UpdateStatus(ctx, gomock.Any(), int64(2)).Return(nil).AnyTimes()

		service.SyncCron(ctx)
	})

	// Test error case
	t.Run("error getting cloud accounts", func(t *testing.T) {
		mHTTP.EXPECT().GetAllCloudAccounts(ctx).
			Return(nil, assert.AnError)
		mocks.Metrics.EXPECT().IncrementCounter(ctx, "sync_error_count")

		service.SyncCron(ctx)
	})
}
