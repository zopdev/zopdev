package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"golang.org/x/oauth2/google"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

func TestService_SyncCron(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cnt, mocks := container.NewMockContainer(t)

	mGCP := NewMockGCPClient(ctrl)
	mHTTP := NewMockHTTPClient(ctrl)
	mStore := NewMockStore(ctrl)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: cnt,
	}
	mockResp := []models.Resource{
		{Name: "sql-instance-1", UID: "zop/sql1"}, {Name: "sql-instance-2", UID: "zop/sql2"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}
	mockCreds := &google.Credentials{ProjectID: "test-project"}

	service := New(mGCP, mHTTP, mStore)

	// mock expectations
	mHTTP.EXPECT().GetAllCloudAccounts(ctx).
		Return([]client.CloudAccount{{ID: 1}, {ID: 2}}, nil)
	mHTTP.EXPECT().GetCloudCredentials(ctx, int64(1)).
		Return(&client.CloudAccount{ID: 1, Provider: "GCP"}, nil)
	mHTTP.EXPECT().GetCloudCredentials(ctx, int64(2)).
		Return(&client.CloudAccount{ID: 2, Provider: "Unknown"}, nil)

	mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
		Return(mockCreds, nil)
	mGCP.EXPECT().NewSQLClient(ctx, gomock.Any()).
		Return(mockLister, nil)

	mStore.EXPECT().GetResources(ctx, int64(1), nil).
		Return(mockResp, nil).Times(2)
	mStore.EXPECT().UpdateResource(ctx, gomock.Any()).
		Return(nil).Times(2)

	// This means that the syncAll returned an error.
	mocks.Metrics.EXPECT().IncrementCounter(ctx, "sync_error_count")

	service.SyncCron(ctx)

	// Failure case: if the HTTP client fails to get cloud accounts
	mHTTP.EXPECT().GetAllCloudAccounts(ctx).
		Return(nil, assert.AnError)
	mocks.Metrics.EXPECT().IncrementCounter(ctx, "sync_error_count")

	service.SyncCron(ctx)
}
