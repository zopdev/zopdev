package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"

	"gofr.dev/pkg/gofr"

	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zopdev/zopdev/api/audit/store"
)

var (
	errInvalidGCPCreds        = errors.New("invalid GCP credentials")
	errInvalidJSONCredentials = errors.New("invalid JSON credentials")
	errCreateSQLAdminService  = errors.New("failed to create SQL Admin service")
	errListCloudSQLInstances  = errors.New("failed to list CloudSQL instances")
	errCreateMonitoringClient = errors.New("failed to create Monitoring client")
	errReadingTimeSeries      = errors.New("error reading time series for sql instance")
)

const (
	// Status levels used to classify CPU utilization.
	danger    = "danger"    // CPU usage is too high or too low; system may be under-provisioned or overloaded.
	warning   = "warning"   // CPU usage is  within tolerable limits.
	compliant = "compliant" // CPU usage is within expected operational range.

	// CPU utilization thresholds (in percentage).
	lowerBound   = 20  // Below this is considered non-compliant due to over-provisioning (under-utilized).
	warningBound = 70  // Between lowerBound and warningBound is considered compliant.
	upperBound   = 90  // Between warningBound and upperBound is warning; above this is danger.
	percentage   = 100 // Represents the full scale (100%) of CPU usage.
)

// CheckCloudSQLProvisionedUsage checks the provisioned usage of Cloud SQL instances
// in a given Google Cloud project. It retrieves the list of Cloud SQL instances
// and their utilization metrics using the Google Cloud SQL Admin API and the
// Cloud Monitoring API.
func CheckCloudSQLProvisionedUsage(ctx *gofr.Context, creds any) ([]store.Items, error) {
	cred, err := getGoogleCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}

	sqlService, err := sqladmin.NewService(ctx, option.WithCredentials(cred))
	if err != nil {
		ctx.Errorf("failed to create SQL Admin service: %v", err)
		return nil, errCreateSQLAdminService
	}

	instancesList, err := sqlService.Instances.List(cred.ProjectID).Do()
	if err != nil {
		ctx.Errorf("failed to list instances: %v", err)
		return nil, errListCloudSQLInstances
	}

	monitoringClient, err := monitoring.NewMetricClient(ctx, option.WithCredentials(cred))
	if err != nil {
		ctx.Errorf("failed to create monitoring client: %v", err)
		return nil, errCreateMonitoringClient
	}

	defer monitoringClient.Close()

	return getResult(ctx, cred.ProjectID, instancesList, monitoringClient)
}

func getResult(ctx *gofr.Context, projectID string,
	instancesList *sqladmin.InstancesListResponse, monitoringClient *monitoring.MetricClient) ([]store.Items, error) {
	results := make([]store.Items, 0)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour) // Token last 24 hours to avergage out the utilization to avoid any outliers
	mu, errGrp := sync.Mutex{}, new(errgroup.Group)

	for _, instance := range instancesList.Items {
		errGrp.Go(func() error {
			resourceFilter := fmt.Sprintf(`resource.type="cloudsql_database" AND resource.labels.database_id=%q`, projectID+":"+instance.Name)
			req := &monitoringpb.ListTimeSeriesRequest{
				Name:   "projects/" + projectID,
				Filter: `metric.type="cloudsql.googleapis.com/database/cpu/utilization" AND ` + resourceFilter,
				Interval: &monitoringpb.TimeInterval{
					StartTime: timestamppb.New(startTime),
					EndTime:   timestamppb.New(endTime),
				},
				View: monitoringpb.ListTimeSeriesRequest_FULL,
			}
			it := monitoringClient.ListTimeSeries(ctx, req)

			var peakUsage float64

			for {
				resp, er := it.Next()
				if errors.Is(er, iterator.Done) {
					break
				}

				if er != nil {
					ctx.Errorf("error reading time series: %v, intance: %s", er, instance.Name)
					return errReadingTimeSeries
				}

				points := resp.Points
				if len(points) > 0 {
					for _, point := range points {
						peakUsage = max(peakUsage, point.Value.GetDoubleValue()*percentage)
					}
				}
			}

			status := compliant
			if peakUsage <= lowerBound {
				status = danger
			}

			meta := map[string]any{
				"peak_utilization": peakUsage,
			}

			mu.Lock()
			results = append(results, store.Items{
				InstanceName: instance.Name,
				Status:       status,
				Metadata:     meta,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := errGrp.Wait(); err != nil {
		return results, err
	}

	return results, nil
}

func getGoogleCredentials(ctx context.Context, creds any) (*google.Credentials, error) {
	if creds == nil {
		return nil, errInvalidGCPCreds
	}

	b, err := json.Marshal(creds)
	if err != nil {
		return nil, errInvalidGCPCreds
	}

	var gcpCred Credentials

	err = json.Unmarshal(b, &gcpCred)
	if err != nil {
		return nil, errInvalidGCPCreds
	}

	cred, err := google.CredentialsFromJSON(ctx, b, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, errInvalidJSONCredentials
	}

	return cred, nil
}
