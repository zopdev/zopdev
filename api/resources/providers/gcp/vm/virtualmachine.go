package vm

import (
	"strings"

	"github.com/zopdev/zopdev/api/resources/providers/models"
	"google.golang.org/api/compute/v1"

	"gofr.dev/pkg/gofr"
)

type ComputeClient struct {
	ComputeService *compute.Service
}

func (c *ComputeClient) GetAllInstances(_ *gofr.Context, projectID string) ([]models.Instance, error) {
	var instances []models.Instance

	aggList, err := c.ComputeService.Instances.AggregatedList(projectID).Do()
	if err != nil {
		return nil, err
	}

	for _, scopedList := range aggList.Items {
		for _, item := range scopedList.Instances {
			if isGKEManaged(item) {
				continue
			}

			zoneParts := strings.Split(item.Zone, "/")
			zone := zoneParts[len(zoneParts)-1]

			instances = append(instances, models.Instance{
				Name:         item.Name,
				Type:         "VM",
				ProviderID:   projectID,
				Region:       zone,
				CreationTime: item.CreationTimestamp,
				Status:       item.Status,
			})
		}
	}

	return instances, nil
}

func isGKEManaged(instance *compute.Instance) bool {
	return hasGKENameOrLabels(instance) || hasGKECreatedByMetadata(instance.Metadata)
}

func hasGKENameOrLabels(instance *compute.Instance) bool {
	if strings.HasPrefix(instance.Name, "gke-") {
		return true
	}

	for key := range instance.Labels {
		if key == "goog-gke-nodepool" || key == "k8s-io-cluster" || key == "gke.io/cluster-name" {
			return true
		}
	}

	return false
}

func hasGKECreatedByMetadata(metadata *compute.Metadata) bool {
	if metadata == nil {
		return false
	}

	for _, item := range metadata.Items {
		if item.Key == "created-by" && item.Value != nil &&
			strings.Contains(*item.Value, "instanceGroupManagers") &&
			strings.Contains(*item.Value, "gke-") {
			return true
		}
	}

	return false
}
