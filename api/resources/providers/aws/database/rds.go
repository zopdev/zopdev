package database

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
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

func (c *Client) GetAllInstances(ctx *gofr.Context) ([]models.Instance, error) {
	input := &rds.DescribeDBInstancesInput{}

	result, err := c.RDS.DescribeDBInstancesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	instances := make([]models.Instance, 0, len(result.DBInstances))

	for _, db := range result.DBInstances {
		engine := awsStringValue(db.Engine)
		clusterID := awsStringValue(db.DBClusterIdentifier)

		var typ string

		switch {
		case strings.Contains(strings.ToLower(engine), "aurora") && clusterID != "":
			typ = "RDS-AURORA"
		case strings.Contains(strings.ToLower(engine), "mysql"):
			typ = "RDS-MYSQL"
		case strings.Contains(strings.ToLower(engine), "postgres"):
			typ = "RDS-POSTGRESQL"
		case strings.Contains(strings.ToLower(engine), "mariadb"):
			typ = "RDS-MARIADB"
		case strings.Contains(strings.ToLower(engine), "oracle"):
			typ = "RDS-ORACLE"
		case strings.Contains(strings.ToLower(engine), "sqlserver"):
			typ = "RDS-SQLSERVER"
		default:
			typ = "RDS-UNKNOWN"
		}

		instance := models.Instance{
			Name:         awsStringValue(db.DBInstanceIdentifier),
			Type:         typ,
			UID:          awsStringValue(db.DBInstanceArn),
			Region:       awsStringValue(db.AvailabilityZone),
			CreationTime: db.InstanceCreateTime.String(),
			Status:       awsStringValue(db.DBInstanceStatus),
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
func (c *Client) StartInstance(ctx *gofr.Context, engine, clusterID, dbInstanceIdentifier string) error {
	engine = strings.ToLower(engine)
	if strings.Contains(engine, "aurora") && clusterID != "" {
		input := &rds.StartDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		}
		_, err := c.RDS.StartDBClusterWithContext(ctx, input)

		return err
	}

	input := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	}
	_, err := c.RDS.StartDBInstanceWithContext(ctx, input)

	return err
}

// StopInstance handles all RDS types: Aurora clusters and standard RDS. Aurora Serverless detection is not supported here.
func (c *Client) StopInstance(ctx *gofr.Context, engine, clusterID, dbInstanceIdentifier string) error {
	engine = strings.ToLower(engine)
	if strings.Contains(engine, "aurora") && clusterID != "" {
		input := &rds.StopDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		}
		_, err := c.RDS.StopDBClusterWithContext(ctx, input)

		return err
	}

	input := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	}
	_, err := c.RDS.StopDBInstanceWithContext(ctx, input)

	return err
}

func awsStringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
