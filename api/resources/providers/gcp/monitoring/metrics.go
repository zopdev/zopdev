package monitoring

import (
	"errors"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Client struct {
	MetricClient *monitoring.MetricClient
}

func (c *Client) GetTimeSeries(ctx *gofr.Context, start, end time.Time, projectID, filter string) ([]models.Metric, error) {
	var (
		metrics = make([]models.Metric, 0)
		req     = &monitoringpb.ListTimeSeriesRequest{
			Name:   "projects/" + projectID,
			Filter: filter,
			Interval: &monitoringpb.TimeInterval{
				StartTime: timestamppb.New(start),
				EndTime:   timestamppb.New(end),
			},
			View: monitoringpb.ListTimeSeriesRequest_FULL,
		}

		// Create an iterator to list time series
		it = c.MetricClient.ListTimeSeries(ctx, req)
	)

	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(resp.Points) == 0 {
			continue
		}

		metrics = append(metrics, models.Metric{
			Point: getTypedValue(resp.Points[0].Value),
		})
	}

	return metrics, nil
}

func getTypedValue(val *monitoringpb.TypedValue) any {
	switch val.Value.(type) {
	case *monitoringpb.TypedValue_BoolValue:
		return val.GetBoolValue()
	case *monitoringpb.TypedValue_DoubleValue:
		return val.GetDoubleValue()
	case *monitoringpb.TypedValue_Int64Value:
		return val.GetInt64Value()
	case *monitoringpb.TypedValue_StringValue:
		return val.GetStringValue()
	default:
		return nil
	}
}
