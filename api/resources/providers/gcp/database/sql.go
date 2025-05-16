package sql

import (
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/sqladmin/v1"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Client struct {
	SQL *sqladmin.InstancesService
}

func (c *Client) GetAllInstances(_ *gofr.Context, projectID string) ([]models.SQLInstance, error) {
	list, err := c.SQL.List(projectID).Do()
	if err != nil {
		return nil, err
	}

	var instances = make([]models.SQLInstance, 0)

	for _, item := range list.Items {
		instances = append(instances, models.SQLInstance{
			Name:         item.Name,
			ProjectID:    item.Project,
			Region:       item.Region,
			Zone:         item.GceZone,
			Version:      item.DatabaseVersion,
			CreationTime: item.CreateTime,
		})
	}

	return instances, nil
}

func (c *Client) StopInstance(_ *gofr.Context, projectID, instanceID string) error {
	_, err := c.SQL.StopReplica(projectID, instanceID).Do()
	if err != nil {
		return err
	}

	return nil
}
