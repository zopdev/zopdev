package sql

import (
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/sqladmin/v1"

	"github.com/zopdev/zopdev/api/resources/models"
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
			InstanceName: item.Name,
			ProjectID:    item.Project,
			Region:       item.Region,
			Zone:         item.GceZone,
			Version:      item.DatabaseVersion,
			CreationTime: item.CreateTime,
		})
	}

	return instances, nil
}
