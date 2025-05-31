package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"golang.org/x/oauth2/google"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/providers/aws/database"
	"github.com/zopdev/zopdev/api/resources/providers/aws/vm"
)

// stubEC2 implements the EC2API interface with no-op methods.
type stubEC2 struct{}

func (*stubEC2) DescribeInstancesWithContext(_ aws.Context, _ *ec2.DescribeInstancesInput,
	_ ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{}, nil
}
func (*stubEC2) StartInstancesWithContext(_ aws.Context, _ *ec2.StartInstancesInput,
	_ ...request.Option) (*ec2.StartInstancesOutput, error) {
	return &ec2.StartInstancesOutput{}, nil
}
func (*stubEC2) StopInstancesWithContext(_ aws.Context, _ *ec2.StopInstancesInput,
	_ ...request.Option) (*ec2.StopInstancesOutput, error) {
	return &ec2.StopInstancesOutput{}, nil
}

// stubRDS implements the RDSAPI interface with no-op methods.
type stubRDS struct{}

func (*stubRDS) DescribeDBInstancesWithContext(_ aws.Context, _ *rds.DescribeDBInstancesInput,
	_ ...request.Option) (*rds.DescribeDBInstancesOutput, error) {
	return &rds.DescribeDBInstancesOutput{}, nil
}
func (*stubRDS) StartDBClusterWithContext(_ aws.Context, _ *rds.StartDBClusterInput,
	_ ...request.Option) (*rds.StartDBClusterOutput, error) {
	return &rds.StartDBClusterOutput{}, nil
}
func (*stubRDS) StartDBInstanceWithContext(_ aws.Context, _ *rds.StartDBInstanceInput,
	_ ...request.Option) (*rds.StartDBInstanceOutput, error) {
	return &rds.StartDBInstanceOutput{}, nil
}
func (*stubRDS) StopDBClusterWithContext(_ aws.Context, _ *rds.StopDBClusterInput,
	_ ...request.Option) (*rds.StopDBClusterOutput, error) {
	return &rds.StopDBClusterOutput{}, nil
}
func (*stubRDS) StopDBInstanceWithContext(_ aws.Context, _ *rds.StopDBInstanceInput,
	_ ...request.Option) (*rds.StopDBInstanceOutput, error) {
	return &rds.StopDBInstanceOutput{}, nil
}

func TestService_SyncCron(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cnt, mocks := container.NewMockContainer(t)

	mGCP := NewMockGCPClient(ctrl)
	mHTTP := NewMockHTTPClient(ctrl)
	mStore := NewMockStore(ctrl)
	mAWS := NewMockAWSClient(ctrl)
	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: cnt,
	}
	mockResp := []models.Instance{
		{ID: 1, Name: "sql-instance-1", UID: "zop/sql1", Status: RUNNING}, {ID: 2, Name: "sql-instance-2", UID: "zop/sql2", Status: RUNNING},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}
	mockCreds := &google.Credentials{ProjectID: "test-project"}

	service := New(mGCP, mAWS, mHTTP, mStore)

	// mock expectations
	mHTTP.EXPECT().GetAllCloudAccounts(ctx).
		Return([]client.CloudAccount{{ID: 1, Provider: "GCP"}, {ID: 2, Provider: "Unknown"}}, nil)
	mHTTP.EXPECT().GetCloudCredentials(ctx, int64(1)).
		Return(&client.CloudAccount{ID: 1, Provider: "GCP"}, nil)
	mHTTP.EXPECT().GetCloudCredentials(ctx, int64(2)).
		Return(&client.CloudAccount{ID: 2, Provider: "Unknown"}, nil)

	mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
		Return(mockCreds, nil)
	mGCP.EXPECT().NewSQLClient(ctx, gomock.Any()).
		Return(mockLister, nil)

	mStore.EXPECT().GetResources(ctx, int64(1), nil).
		Return(mockResp, nil).AnyTimes()
	mStore.EXPECT().GetResources(ctx, int64(2), nil).
		Return(nil, nil).AnyTimes()
	mStore.EXPECT().UpdateStatus(ctx, RUNNING, int64(1)).
		Return(nil)
	mStore.EXPECT().UpdateStatus(ctx, RUNNING, int64(2)).
		Return(nil)

	// Add correct mocks for AWS EC2 and RDS clients
	mAWS.EXPECT().NewEC2Client(gomock.Any(), gomock.Any()).Return(&vm.Client{EC2: &stubEC2{}}, nil).AnyTimes()
	mAWS.EXPECT().NewRDSClient(gomock.Any(), gomock.Any()).Return(&database.Client{RDS: &stubRDS{}}, nil).AnyTimes()

	service.SyncCron(ctx)

	// Failure case: if the HTTP client fails to get cloud accounts
	mHTTP.EXPECT().GetAllCloudAccounts(ctx).
		Return(nil, assert.AnError)
	mocks.Metrics.EXPECT().IncrementCounter(ctx, "sync_error_count")

	service.SyncCron(ctx)
}
