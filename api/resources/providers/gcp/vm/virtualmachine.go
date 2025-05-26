package vm

import (
	"github.com/zopdev/zopdev/api/resources/providers/models"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/compute/v1"
	"strings"
)

type ComputeClient struct {
	ComputeService *compute.Service
}

func (c *ComputeClient) GetAllVMInstances(ctx *gofr.Context, projectID string) ([]models.Instance, error) {
	var instances []models.Instance

	aggList, err := c.ComputeService.Instances.AggregatedList(projectID).Do()
	if err != nil {
		return nil, err
	}

	for _, scopedList := range aggList.Items {
		for _, item := range scopedList.Instances {
			zoneParts := strings.Split(item.Zone, "/")
			zone := zoneParts[len(zoneParts)-1]
			region := zone[:len(zone)-2]

			instances = append(instances, models.Instance{
				Name:         item.Name,
				Type:         "VM",
				ProviderID:   projectID,
				Region:       region,
				CreationTime: item.CreationTimestamp,
				Status:       item.Status,
			})
		}
	}

	return instances, nil
}

func (c *ComputeClient) StartInstanceVM(_ *gofr.Context, projectID, zone, instanceName string) error {
	_, err := c.ComputeService.Instances.Start(projectID, zone, instanceName).Do()
	if err != nil {
		return err
	}
	return nil
}

func (c *ComputeClient) StopInstanceVM(_ *gofr.Context, projectID, zone, instanceName string) error {
	_, err := c.ComputeService.Instances.Stop(projectID, zone, instanceName).Do()
	if err != nil {
		return err
	}
	return nil
}
