package database

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRDS struct {
	dbInstances []*rds.DBInstance
	shouldErr   bool
}

func (m *mockRDS) DescribeDBInstancesWithContext(_ aws.Context, _ *rds.DescribeDBInstancesInput,
	_ ...request.Option) (*rds.DescribeDBInstancesOutput, error) {
	if m.shouldErr {
		return nil, assert.AnError
	}

	return &rds.DescribeDBInstancesOutput{DBInstances: m.dbInstances}, nil
}

func (m *mockRDS) StartDBClusterWithContext(_ aws.Context, _ *rds.StartDBClusterInput,
	_ ...request.Option) (*rds.StartDBClusterOutput, error) {
	if m.shouldErr {
		return nil, assert.AnError
	}

	return &rds.StartDBClusterOutput{}, nil
}

func (m *mockRDS) StartDBInstanceWithContext(_ aws.Context, _ *rds.StartDBInstanceInput,
	_ ...request.Option) (*rds.StartDBInstanceOutput, error) {
	if m.shouldErr {
		return nil, assert.AnError
	}

	return &rds.StartDBInstanceOutput{}, nil
}

func (m *mockRDS) StopDBClusterWithContext(_ aws.Context, _ *rds.StopDBClusterInput,
	_ ...request.Option) (*rds.StopDBClusterOutput, error) {
	if m.shouldErr {
		return nil, assert.AnError
	}

	return &rds.StopDBClusterOutput{}, nil
}

func (m *mockRDS) StopDBInstanceWithContext(_ aws.Context, _ *rds.StopDBInstanceInput,
	_ ...request.Option) (*rds.StopDBInstanceOutput, error) {
	if m.shouldErr {
		return nil, assert.AnError
	}

	return &rds.StopDBInstanceOutput{}, nil
}

func Test_GetAllInstances(t *testing.T) {
	mock := &mockRDS{
		dbInstances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test-rds-1"),
				DBInstanceArn:        aws.String("arn:aws:rds:us-east-1:123456789012:db:test-rds-1"),
				AvailabilityZone:     aws.String("us-east-1a"),
				InstanceCreateTime:   aws.Time(time.Now()),
				DBInstanceStatus:     aws.String("available"),
				Engine:               aws.String("mysql"),
			},
		},
	}
	client := &Client{RDS: mock}
	instances, err := client.GetAllInstances(nil)
	require.NoError(t, err)
	require.Len(t, instances, 1)
	assert.Equal(t, "test-rds-1", instances[0].Name)
	assert.Equal(t, "RDS-MYSQL", instances[0].Type)
	assert.Equal(t, "us-east-1a", instances[0].Region)
	assert.Equal(t, "available", instances[0].Status)
}

func Test_GetAllInstances_Error(t *testing.T) {
	mock := &mockRDS{shouldErr: true}
	client := &Client{RDS: mock}
	instances, err := client.GetAllInstances(nil)
	require.Error(t, err)
	require.Nil(t, instances)
}

func Test_StartInstance(t *testing.T) {
	cases := []struct {
		name      string
		engine    string
		clusterID string
		shouldErr bool
	}{
		{"Aurora Success", "aurora-mysql", "test-cluster", false},
		{"Standard Success", "mysql", "", false},
		{"Aurora Error", "aurora-mysql", "test-cluster", true},
		{"Standard Error", "mysql", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mockRDS{shouldErr: c.shouldErr}
			client := &Client{RDS: mock}

			err := client.StartInstance(nil, c.engine, c.clusterID, "test-instance")
			if c.shouldErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StopInstance(t *testing.T) {
	cases := []struct {
		name      string
		engine    string
		clusterID string
		shouldErr bool
	}{
		{"Aurora Success", "aurora-mysql", "test-cluster", false},
		{"Standard Success", "mysql", "", false},
		{"Aurora Error", "aurora-mysql", "test-cluster", true},
		{"Standard Error", "mysql", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mockRDS{shouldErr: c.shouldErr}
			client := &Client{RDS: mock}

			err := client.StopInstance(nil, c.engine, c.clusterID, "test-instance")
			if c.shouldErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GetAllInstances_AllTypes(t *testing.T) {
	engines := []struct {
		engine    string
		clusterID string
		expType   string
	}{
		{"aurora-mysql", "cluster-1", "RDS-AURORA"},
		{"mysql", "", "RDS-MYSQL"},
		{"postgres", "", "RDS-POSTGRESQL"},
		{"mariadb", "", "RDS-MARIADB"},
		{"oracle", "", "RDS-ORACLE"},
		{"sqlserver", "", "RDS-SQLSERVER"},
		{"customengine", "", "RDS-UNKNOWN"},
	}
	for _, e := range engines {
		db := &rds.DBInstance{
			DBInstanceIdentifier: aws.String("id-" + e.expType),
			DBInstanceArn:        aws.String("arn:aws:rds:region:acct:db:id-" + e.expType),
			AvailabilityZone:     aws.String("region-1a"),
			InstanceCreateTime:   aws.Time(time.Now()),
			DBInstanceStatus:     aws.String("available"),
			Engine:               aws.String(e.engine),
			DBClusterIdentifier:  aws.String(e.clusterID),
		}
		mock := &mockRDS{dbInstances: []*rds.DBInstance{db}}
		client := &Client{RDS: mock}
		instances, err := client.GetAllInstances(nil)
		require.NoError(t, err)
		require.Len(t, instances, 1)
		assert.Equal(t, e.expType, instances[0].Type)
	}
}
