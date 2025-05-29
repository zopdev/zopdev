package vm

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zopdev/zopdev/api/resources/providers/models"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getComputeServer(t *testing.T, response any) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if response == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "unable to marshal response", http.StatusInternalServerError)
			return
		}
	}))

	return srv
}

func TestGetAllVMInstances(t *testing.T) {
	projectID := "test-project"

	mockResponse := getMockAggregatedInstanceList()
	expected := getExpectedVMInstances(projectID)

	server := getComputeServer(t, mockResponse)
	defer server.Close()

	computeSvc, err := compute.NewService(context.TODO(), option.WithoutAuthentication(), option.WithEndpoint(server.URL))
	require.NoError(t, err)

	client := ComputeClient{ComputeService: computeSvc}

	instances, err := client.GetAllInstances(nil, projectID)
	require.NoError(t, err)

	assert.Equal(t, expected, instances)
}

func getMockAggregatedInstanceList() *compute.InstanceAggregatedList {
	return &compute.InstanceAggregatedList{
		Items: map[string]compute.InstancesScopedList{
			"zones/us-central1": {
				Instances: getMockInstances(),
			},
		},
	}
}

func getMockInstances() []*compute.Instance {
	gkeCreatedByValue := "projects/test-project/zones/us-central1-a/instanceGroupManagers/gke-nodepool"
	nonGKECreatedByValue := "projects/test-project/zones/us-central1-a/instanceGroupManagers/custom-instance-group"

	return []*compute.Instance{
		{
			Name:              "normal-vm",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
			CreationTimestamp: "2025-01-01T00:00:00.000-07:00",
			Status:            "RUNNING",
		},
		{
			Name:              "managed-vm",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
			CreationTimestamp: "2025-01-02T00:00:00.000-07:00",
			Status:            "RUNNING",
			Metadata: &compute.Metadata{
				Items: []*compute.MetadataItems{
					{
						Key:   "created-by",
						Value: &gkeCreatedByValue,
					},
				},
			},
		},
		{
			Name:              "non-gke-vm-with-created-by",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
			CreationTimestamp: "2025-01-03T00:00:00.000-07:00",
			Status:            "RUNNING",
		},
		{
			Name:              "vm-with-gke-io-label",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
			CreationTimestamp: "2025-07-01T00:00:00.000-07:00",
			Status:            "RUNNING",
			Labels: map[string]string{
				"gke.io/cluster-name": "cluster-example",
			},
		},
		{
			Name:              "gke-vm",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1",
			CreationTimestamp: "2025-03-01T00:00:00.000-07:00",
			Status:            "RUNNING",
		},
		{
			Name:              "vm-with-non-gke-created-by",
			Zone:              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
			CreationTimestamp: "2025-07-05T00:00:00.000-07:00",
			Status:            "RUNNING",
			Metadata: &compute.Metadata{
				Items: []*compute.MetadataItems{
					{
						Key:   "created-by",
						Value: &nonGKECreatedByValue,
					},
				},
			},
		},
	}
}

func getExpectedVMInstances(projectID string) []models.Instance {
	return []models.Instance{
		{
			Name:         "normal-vm",
			Type:         "VM",
			ProviderID:   projectID,
			Region:       "us-central1-a",
			CreationTime: "2025-01-01T00:00:00.000-07:00",
			Status:       "RUNNING",
		},
		{
			Name:         "non-gke-vm-with-created-by",
			Type:         "VM",
			ProviderID:   projectID,
			Region:       "us-central1-a",
			CreationTime: "2025-01-03T00:00:00.000-07:00",
			Status:       "RUNNING",
		},
		{
			Name:         "vm-with-non-gke-created-by",
			Type:         "VM",
			ProviderID:   projectID,
			Region:       "us-central1-a",
			CreationTime: "2025-07-05T00:00:00.000-07:00",
			Status:       "RUNNING",
		},
	}
}

func Test_GetAllInstances_Error(t *testing.T) {
	srv := getComputeServer(t, nil)
	defer srv.Close()

	computeService, err := compute.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(srv.URL),
	)
	require.NoError(t, err)

	client := ComputeClient{
		ComputeService: computeService,
	}

	instances, err := client.GetAllInstances(nil, "test-project")
	require.Error(t, err)
	require.Nil(t, instances)

	expected := &googleapi.Error{
		Code: http.StatusInternalServerError,
		Body: "Internal server error\n",
	}
	assert.Equal(t, expected.Error(), err.Error())
}
