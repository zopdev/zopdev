package monitoring

import (
	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

var errMock = errors.New("error")

type fakeMetricServer struct {
	isError bool
	monitoringpb.UnimplementedMetricServiceServer
}

func (f *fakeMetricServer) ListTimeSeries(_ context.Context,
	_ *monitoringpb.ListTimeSeriesRequest) (*monitoringpb.ListTimeSeriesResponse, error) {
	if f.isError {
		return nil, errMock
	}

	resp := &monitoringpb.ListTimeSeriesResponse{
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Points: []*monitoringpb.Point{
					{
						Value: &monitoringpb.TypedValue{
							Value: &monitoringpb.TypedValue_Int64Value{Int64Value: 1},
						},
					},
				},
			},
			{
				Points: []*monitoringpb.Point{
					{
						Value: &monitoringpb.TypedValue{
							Value: &monitoringpb.TypedValue_BoolValue{BoolValue: false},
						},
					},
				},
			},
			{
				Points: []*monitoringpb.Point{
					{
						Value: &monitoringpb.TypedValue{
							Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: 12.3},
						},
					},
				},
			},
			{
				Points: []*monitoringpb.Point{
					{
						Value: &monitoringpb.TypedValue{
							Value: &monitoringpb.TypedValue_StringValue{StringValue: "29"},
						},
					},
				},
			},
			{
				Points: []*monitoringpb.Point{},
			},
		},
	}

	return resp, nil
}

func getGRPCServer(t *testing.T, isError bool) (*grpc.Server, string) {
	t.Helper()

	fakeServer := &fakeMetricServer{isError: isError}

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	grpcSrv := grpc.NewServer()
	fakeServerAddr := l.Addr().String()

	monitoringpb.RegisterMetricServiceServer(grpcSrv, fakeServer)

	go func() {
		if er := grpcSrv.Serve(l); er != nil {
			panic(er)
		}
	}()

	return grpcSrv, fakeServerAddr
}

func TestClient_GetTimeSeries(t *testing.T) {
	grpcSrv, fakeServerAddr := getGRPCServer(t, false)
	defer grpcSrv.Stop()

	metricClient, err := monitoring.NewMetricClient(
		context.Background(),
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	if err != nil {
		t.Fatal(err)
	}

	client := Client{metricClient}
	startTime := time.Now()
	endTime := startTime.Add(-24 * time.Hour)
	expected := []models.Metric{
		{Point: int64(1)},
		{Point: false},
		{Point: 12.3},
		{Point: "29"},
	}

	resp, er := client.GetTimeSeries(&gofr.Context{Context: context.Background()}, startTime, endTime, "test-project", "")

	require.NoError(t, er)
	assert.Equal(t, expected, resp)
}

func TestClient_GetTimeSeries_Error(t *testing.T) {
	grpcSrv, fakeServerAddr := getGRPCServer(t, true)
	defer grpcSrv.Stop()

	metricClient, err := monitoring.NewMetricClient(
		context.Background(),
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	if err != nil {
		t.Fatal(err)
	}

	client := Client{metricClient}
	startTime := time.Now()
	endTime := startTime.Add(-24 * time.Hour)
	resp, er := client.GetTimeSeries(&gofr.Context{Context: context.Background()}, startTime, endTime, "test-project", "")

	require.Nil(t, resp)
	assert.ErrorIs(t, errMock, er)
}
