package oci

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"gofr.dev/pkg/gofr"
	"golang.org/x/sync/errgroup"

	"github.com/zopdev/zopdev/api/audit/store"
)

var (
	errInvalidOCICreds  = errors.New("invalid OCI credentials")
	errCreateDBClient   = errors.New("failed to create Database client")
	errListDBSystems    = errors.New("failed to list DB systems")
	errMonitoringClient = errors.New("failed to create Monitoring client")
	errReadingMetrics   = errors.New("error reading metrics for DB system")
)

const (
	danger    = "danger"
	warning   = "warning"
	compliant = "compliant"

	lowerBound   = 20
	warningBound = 70
	upperBound   = 90
	percentage   = 100

	metricNamespace = "oci_database"
	metricName      = "CpuUtilization"
)

// CheckDBSystemProvisionedUsage checks the provisioned usage of DB systems
// in a given OCI compartment. It retrieves the list of DB systems
// and their utilization metrics using the OCI Database and Monitoring APIs.
func CheckDBSystemProvisionedUsage(ctx *gofr.Context, creds any) ([]store.Items, error) {
	ociCreds, err := getOCICredentials(creds)
	if err != nil {
		return nil, err
	}

	// Format the private key properly by replacing literal \n with actual newlines
	privateKey := strings.ReplaceAll(ociCreds.PrivateKey, "\\n", "\n")

	configProvider := common.NewRawConfigurationProvider(
		ociCreds.TenancyOCID,
		ociCreds.UserOCID,
		ociCreds.Region,
		ociCreds.Fingerprint,
		privateKey,
		nil,
	)

	dbClient, err := database.NewDatabaseClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, errCreateDBClient
	}

	// List DB Systems
	dbSystems, err := listDBSystems(ctx, &dbClient, ociCreds.Compartment)
	if err != nil {
		return nil, err
	}

	monitoringClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, errMonitoringClient
	}

	return getResult(ctx, ociCreds, dbSystems, &monitoringClient)
}

func listDBSystems(ctx context.Context, client *database.DatabaseClient, compartmentID string) ([]database.DbSystemSummary, error) {
	request := database.ListDbSystemsRequest{
		CompartmentId: &compartmentID,
	}

	response, err := client.ListDbSystems(ctx, request)
	if err != nil {
		return nil, errListDBSystems
	}

	return response.Items, nil
}

func getResult(ctx context.Context, creds *Credentials, dbSystems []database.DbSystemSummary,
	monitoringClient *monitoring.MonitoringClient) ([]store.Items, error) {
	results := make([]store.Items, 0)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	mu, errGrp := sync.Mutex{}, new(errgroup.Group)

	for i := range dbSystems {
		system := &dbSystems[i]

		errGrp.Go(func() error {
			request := monitoring.SummarizeMetricsDataRequest{
				CompartmentId: &creds.Compartment,
				SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
					Namespace: common.String(metricNamespace),
					Query:     common.String(fmt.Sprintf("%s[1d]{resourceId = %q}.max()", metricName, *system.Id)),
					StartTime: &common.SDKTime{Time: startTime},
					EndTime:   &common.SDKTime{Time: endTime},
				},
			}

			response, err := monitoringClient.SummarizeMetricsData(ctx, request)
			if err != nil {
				return errReadingMetrics
			}

			var peakUsage float64
			if len(response.Items) > 0 && len(response.Items[0].AggregatedDatapoints) > 0 {
				peakUsage = *response.Items[0].AggregatedDatapoints[0].Value * percentage
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
				InstanceName: *system.DisplayName,
				Status:       status,
				Metadata:     meta,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := errGrp.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func getOCICredentials(creds any) (*Credentials, error) {
	if creds == nil {
		return nil, errInvalidOCICreds
	}

	b, err := json.Marshal(creds)
	if err != nil {
		return nil, errInvalidOCICreds
	}

	var ociCred Credentials
	if err := json.Unmarshal(b, &ociCred); err != nil {
		return nil, errInvalidOCICreds
	}

	return &ociCred, nil
}
