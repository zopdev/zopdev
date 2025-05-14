package monitoring

import (
	"errors"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zopdev/zopdev/api/resources/models"
)

type Client struct {
	MetricClient *monitoring.MetricClient
}

func (c *Client) GetTimeSeries(ctx *gofr.Context, start, end time.Time, projectID, filter string) ([]models.Metric, error) {
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   "projects/" + projectID,
		Filter: filter,
		Interval: &monitoringpb.TimeInterval{
			StartTime: timestamppb.New(start),
			EndTime:   timestamppb.New(end),
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	it := c.MetricClient.ListTimeSeries(ctx, req)
	var metrics []models.Metric

	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, err
		}

		metrics = append(metrics, models.Metric{
			Points: getPoints(resp.GetPoints()),
		})
	}

	return metrics, nil
}

func getPoints(points []*monitoringpb.Point) []models.Point {
	var p []models.Point

	for _, point := range points {
		p = append(p, models.Point{
			Value: getTypedValue(point.Value),
		})
	}

	return p
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
