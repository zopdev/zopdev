package database

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
	gofrService "gofr.dev/pkg/gofr/http"
)

const (
	// RUNNING instance state for zopdev.
	RUNNING = "RUNNING"
	// STOPPED instance state for zopdev.
	STOPPED = "STOPPED"
)

// RDSAPI defines the methods used from the AWS RDS client for easier testing/mocking.
type RDSAPI interface {
	DescribeDBInstancesWithContext(ctx aws.Context, input *rds.DescribeDBInstancesInput,
		opts ...request.Option) (*rds.DescribeDBInstancesOutput, error)
	StartDBClusterWithContext(ctx aws.Context, input *rds.StartDBClusterInput, opts ...request.Option) (*rds.StartDBClusterOutput, error)
	StartDBInstanceWithContext(ctx aws.Context, input *rds.StartDBInstanceInput, opts ...request.Option) (*rds.StartDBInstanceOutput, error)
	StopDBClusterWithContext(ctx aws.Context, input *rds.StopDBClusterInput, opts ...request.Option) (*rds.StopDBClusterOutput, error)
	StopDBInstanceWithContext(ctx aws.Context, input *rds.StopDBInstanceInput, opts ...request.Option) (*rds.StopDBInstanceOutput, error)
}

type Client struct {
	RDS RDSAPI
}

// mapRDSStatus maps AWS RDS DBInstanceStatus to RUNNING, STOPPED, or the original status.
func mapRDSStatus(status string) string {
	switch strings.ToLower(status) {
	case "available", "backing-up", "configuring-enhanced-monitoring", "configuring-iam-database-auth",
		"configuring-log-exports", "converting-to-vpc", "maintenance", "modifying", "moving-to-vpc",
		"rebooting", "resetting-master-credentials", "renaming", "restore-error", "storage-config-upgrade",
		"storage-full", "storage-initialization", "storage-optimization", "upgrading":
		return RUNNING
	case "stopped", "stopping", "starting":
		return STOPPED
	default:
		return status
	}
}

// extractEngineAndClusterID extracts and validates engine and cluster_id from resource.Settings.
func extractEngineAndClusterID(resource *models.Resource) (engine, clusterID string, err error) {
	eng, ok := resource.Settings["engine"].(string)
	if !ok {
		return "", "", gofrService.ErrorInvalidParam{Params: []string{"resource.Settings.engine"}}
	}

	clID, ok := resource.Settings["cluster_id"].(string)
	if !ok {
		return "", "", gofrService.ErrorInvalidParam{Params: []string{"resource.Settings.cluster_id"}}
	}

	return eng, clID, nil
}

func (c *Client) GetAllInstances(ctx *gofr.Context) ([]models.Resource, error) {
	input := &rds.DescribeDBInstancesInput{}

	result, err := c.RDS.DescribeDBInstancesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	instances := make([]models.Resource, 0, len(result.DBInstances))

	for _, db := range result.DBInstances {
		engine := awsStringValue(db.Engine)
		clusterID := awsStringValue(db.DBClusterIdentifier)
		status := awsStringValue(db.DBInstanceStatus)

		mappedStatus := mapRDSStatus(status)

		instance := models.Resource{
			Name:         awsStringValue(db.DBInstanceIdentifier),
			Type:         "RDS",
			UID:          awsStringValue(db.DBInstanceArn),
			Region:       awsStringValue(db.AvailabilityZone),
			CreationTime: db.InstanceCreateTime.String(),
			Status:       mappedStatus,
			CloudAccount: models.CloudAccount{}, // TODO: Set from context or parameter if available
			Settings: map[string]any{
				"engine":     engine,
				"cluster_id": clusterID,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// StartInstance handles all RDS types: Aurora clusters and standard RDS. Aurora Serverless detection is not supported here.
func (c *Client) StartInstance(ctx *gofr.Context, resource *models.Resource) error {
	engine, clusterID, err := extractEngineAndClusterID(resource)
	if err != nil {
		return err
	}

	engine = strings.ToLower(engine)

	if strings.Contains(engine, "aurora") && clusterID != "" {
		input := &rds.StartDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		}
		_, err = c.RDS.StartDBClusterWithContext(ctx, input)

		return err
	}

	input := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: aws.String(resource.Name),
	}
	_, err = c.RDS.StartDBInstanceWithContext(ctx, input)

	return err
}

// StopInstance handles all RDS types: Aurora clusters and standard RDS. Aurora Serverless detection is not supported here.
func (c *Client) StopInstance(ctx *gofr.Context, resource *models.Resource) error {
	engine, clusterID, err := extractEngineAndClusterID(resource)
	if err != nil {
		return err
	}

	engine = strings.ToLower(engine)

	if strings.Contains(engine, "aurora") && clusterID != "" {
		input := &rds.StopDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		}
		_, err = c.RDS.StopDBClusterWithContext(ctx, input)

		return err
	}

	input := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: aws.String(resource.Name),
	}
	_, err = c.RDS.StopDBInstanceWithContext(ctx, input)

	return err
}

func awsStringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
