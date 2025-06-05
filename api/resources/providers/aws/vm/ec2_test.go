package vm

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errFail = errors.New("fail")

type mockEC2 struct {
	DescribeInstancesResp map[string]*ec2.DescribeInstancesOutput
	DescribeInstancesErr  error
	StartErr              error
	StopErr               error
	Region                string
	CurrentRegion         string
}

func (m *mockEC2) DescribeInstancesWithContext(_ aws.Context, _ *ec2.DescribeInstancesInput,
	_ ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	if m.DescribeInstancesErr != nil {
		return nil, m.DescribeInstancesErr
	}

	resp, ok := m.DescribeInstancesResp[m.CurrentRegion]
	if !ok {
		return &ec2.DescribeInstancesOutput{}, nil
	}

	return resp, nil
}

func (m *mockEC2) StartInstancesWithContext(_ aws.Context, _ *ec2.StartInstancesInput,
	_ ...request.Option) (*ec2.StartInstancesOutput, error) {
	return &ec2.StartInstancesOutput{}, m.StartErr
}

func (m *mockEC2) StopInstancesWithContext(_ aws.Context, _ *ec2.StopInstancesInput,
	_ ...request.Option) (*ec2.StopInstancesOutput, error) {
	return &ec2.StopInstancesOutput{}, m.StopErr
}

func Test_GetAllInstances_Success(t *testing.T) {
	regions := GetAWSRegions()
	mock := &mockEC2{
		DescribeInstancesResp: map[string]*ec2.DescribeInstancesOutput{},
	}

	for _, region := range regions {
		mock.DescribeInstancesResp[region] = &ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{{
				Instances: []*ec2.Instance{{
					InstanceId:   aws.String("i-123"),
					InstanceType: aws.String("t2.micro"),
					LaunchTime:   aws.Time(time.Now()),
					State:        &ec2.InstanceState{Name: aws.String("running")},
					Tags:         []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-instance")}},
				}},
			}},
		}
	}

	mock.CurrentRegion = "us-east-1" // Not used in this test, but set for completeness
	client := &Client{EC2: mock}
	instances, err := client.GetAllInstances(nil)
	require.NoError(t, err)
	require.Len(t, instances, len(regions))

	regionSet := make(map[string]struct{})

	for _, inst := range instances {
		assert.Equal(t, "test-instance", inst.Name)
		assert.Equal(t, "EC2", inst.Type)
		assert.Equal(t, "i-123", inst.UID)
		assert.Equal(t, "running", inst.Status)

		regionSet[inst.Region] = struct{}{}
	}
	// Ensure all regions are present in the results
	for _, region := range regions {
		_, ok := regionSet[region]
		assert.True(t, ok, "region %s not found in results", region)
	}
}

func Test_GetAllInstances_Error(t *testing.T) {
	mock := &mockEC2{DescribeInstancesErr: errFail, Region: "us-east-1"}
	mock.CurrentRegion = "us-east-1"
	client := &Client{EC2: mock}
	instances, err := client.GetAllInstances(nil)
	require.Error(t, err)
	require.Empty(t, instances)
}

func Test_GetAllInstances_NoReservations(t *testing.T) {
	mock := &mockEC2{
		DescribeInstancesResp: map[string]*ec2.DescribeInstancesOutput{
			"us-east-1": {Reservations: []*ec2.Reservation{}},
		},
		Region: "us-east-1",
	}
	mock.CurrentRegion = "us-east-1"
	client := &Client{EC2: mock}
	instances, err := client.GetAllInstances(nil)
	require.NoError(t, err)
	require.Empty(t, instances)
}

func Test_StartInstance(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expectOK bool
	}{
		{"Success", nil, true},
		{"Error", errFail, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mockEC2{StartErr: c.err}
			mock.CurrentRegion = "us-east-1"
			client := &Client{EC2: mock}

			err := client.StartInstance(nil, "i-123")
			if c.expectOK {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_StopInstance(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expectOK bool
	}{
		{"Success", nil, true},
		{"Error", errFail, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mockEC2{StopErr: c.err}
			mock.CurrentRegion = "us-east-1"
			client := &Client{EC2: mock}

			err := client.StopInstance(nil, "i-123")
			if c.expectOK {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_awsStringValue(t *testing.T) {
	var nilStr *string

	require.Empty(t, awsStringValue(nilStr))

	val := "hello"
	assert.Equal(t, "hello", awsStringValue(&val))
}
