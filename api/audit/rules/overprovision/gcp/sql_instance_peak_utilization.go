package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/zopdev/zopdev/api/audit/store"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CheckCloudSQLProvisionedUsage(ctx *gofr.Context, creds *Credentials) ([]store.Items, error) {
	var results []store.Items

	b, _ := json.Marshal(creds)

	cred, err := google.CredentialsFromJSON(ctx, b, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		ctx.Logger.Errorf("invalid json credentials: %v", err)

		return nil, fmt.Errorf("invalid json credentials: %v", err)
	}

	// Create Cloud SQL Admin client
	sqlService, err := sqladmin.NewService(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to create SQL Admin service: %v", err)
	}

	// List instances
	instancesList, err := sqlService.Instances.List(creds.ProjectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list CloudSQL instances: %v", err)
	}

	// Create Monitoring client
	monitoringClient, err := monitoring.NewMetricClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to create Monitoring client: %v", err)
	}
	defer monitoringClient.Close()

	endTime := time.Now()
	startTime := endTime.Add(-5 * time.Minute) // last 5 minutes window

	for _, instance := range instancesList.Items {
		resourceFilter := fmt.Sprintf(`resource.type="cloudsql_database" AND resource.labels.database_id="%s"`, instance.Name)

		req := &monitoringpb.ListTimeSeriesRequest{
			Name:   "projects/" + creds.ProjectID,
			Filter: `metric.type="cloudsql.googleapis.com/database/cpu/utilization" AND ` + resourceFilter,
			Interval: &monitoringpb.TimeInterval{
				StartTime: timestamppb.New(startTime),
				EndTime:   timestamppb.New(endTime),
			},
			View: monitoringpb.ListTimeSeriesRequest_FULL,
		}

		it := monitoringClient.ListTimeSeries(ctx, req)

		var latestVal float64
		for {
			resp, er := it.Next()
			if errors.Is(er, iterator.Done) {
				break
			}
			if er != nil {
				return nil, fmt.Errorf("error reading time series for instance %s: %v", instance.Name, err)
			}
			points := resp.Points
			if len(points) > 0 {
				latestVal = points[0].Value.GetDoubleValue() * 100 // utilization in percentage
			}
		}

		// Determine status based on utilization
		status := "passing"
		switch {
		case latestVal <= 20:
			status = "failing"
		case latestVal > 70 && latestVal <= 90:
			status = "warning"
		case latestVal > 90:
			status = "failing"
		}

		results = append(results, store.Items{
			InstanceName: instance.Name,
			Status:       status,
		})
	}

	return results, nil
}
