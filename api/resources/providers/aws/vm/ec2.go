package vm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
)

// EC2API defines the methods used from the AWS EC2 client for easier testing/mocking.
type EC2API interface {
	DescribeInstancesWithContext(ctx aws.Context, input *ec2.DescribeInstancesInput, opts ...request.Option) (*ec2.DescribeInstancesOutput, error)
	StartInstancesWithContext(ctx aws.Context, input *ec2.StartInstancesInput, opts ...request.Option) (*ec2.StartInstancesOutput, error)
	StopInstancesWithContext(ctx aws.Context, input *ec2.StopInstancesInput, opts ...request.Option) (*ec2.StopInstancesOutput, error)
}

type Client struct {
	EC2 EC2API
}

func (c *Client) GetAllInstances(ctx *gofr.Context) ([]models.Instance, error) {
	regions := GetAWSRegions()

	type result struct {
		instances []models.Instance
		err       error
	}

	resultsCh := make(chan result, len(regions))
	for _, region := range regions {
		currentRegion := region // capture range variable
		go func(region string) {
			input := &ec2.DescribeInstancesInput{}
			ec2Result, err := c.EC2.DescribeInstancesWithContext(ctx, input)
			if err != nil {
				resultsCh <- result{nil, err}
				return
			}
			instances := make([]models.Instance, 0)
			for _, reservation := range ec2Result.Reservations {
				for _, inst := range reservation.Instances {

					var instanceName string
					for _, tag := range inst.Tags {
						if *tag.Key == "Name" {
							instanceName = awsStringValue(tag.Value)
							break
						}
					}

					instance := models.Instance{
						Name:         instanceName,
						Type:         fmt.Sprintf("EC2-%v", awsStringValue(inst.InstanceType)),
						UID:          awsStringValue(inst.InstanceId),
						Region:       region,
						CreationTime: inst.LaunchTime.Format(time.RFC3339),
						Status:       awsStringValue(inst.State.Name),
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}
					instances = append(instances, instance)
				}
			}

			resultsCh <- result{instances, nil}
		}(currentRegion)
	}

	allInstances := make([]models.Instance, 0)
	var firstErr error
	for i := 0; i < len(regions); i++ {
		r := <-resultsCh
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
		allInstances = append(allInstances, r.instances...)
	}
	close(resultsCh)

	return allInstances, firstErr
}

func (c *Client) StartInstance(ctx *gofr.Context, instanceID string) error {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{&instanceID},
	}
	_, err := c.EC2.StartInstancesWithContext(ctx, input)
	return err
}

func (c *Client) StopInstance(ctx *gofr.Context, instanceID string) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{&instanceID},
	}
	_, err := c.EC2.StopInstancesWithContext(ctx, input)
	return err
}

func awsStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
