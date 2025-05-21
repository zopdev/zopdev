package sql

import (
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/sqladmin/v1"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Client struct {
	SQL *sqladmin.InstancesService
}

func (c *Client) GetAllInstances(_ *gofr.Context, projectID string) ([]models.Instance, error) {
	list, err := c.SQL.List(projectID).Do()
	if err != nil {
		return nil, err
	}

	var instances = make([]models.Instance, 0)

	for _, item := range list.Items {
		instances = append(instances, models.Instance{
			Name:         item.Name,
			Type:         "SQL",
			ProviderID:   item.Project,
			Region:       item.Region,
			CreationTime: item.CreateTime,
		})
	}

	return instances, nil
}
