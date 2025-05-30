package sql

import (
	"errors"
	"net/http"

	"gofr.dev/pkg/gofr"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sqladmin/v1"

	"github.com/zopdev/zopdev/api/resources/models"
)

const (
	// RUNNING instance state for zopdev.
	RUNNING = "RUNNING"
	// SUSPENDED instance state for zopdev.
	SUSPENDED = "SUSPENDED"

	// The following constants are to identify and change the state of the instance.

	// ALWAYS - The instance is always running.
	ALWAYS = "ALWAYS"
	// NEVER - The instance is never running.
	NEVER = "NEVER"
)

type Client struct {
	SQL *sqladmin.InstancesService
}

func (c *Client) GetAllInstances(_ *gofr.Context, projectID string) ([]models.Resource, error) {
	list, err := c.SQL.List(projectID).Do()
	if err != nil {
		return nil, err
	}

	var instances = make([]models.Resource, 0)

	for _, item := range list.Items {
		instances = append(instances, models.Resource{
			Name:         item.Name,
			Type:         "SQL",
			Region:       item.Region,
			CreationTime: item.CreateTime,
			UID:          projectID + "/" + item.Name,
			Status:       getState(item.Settings.ActivationPolicy),
		})
	}

	return instances, nil
}

func getState(state string) string {
	switch state {
	case ALWAYS:
		return RUNNING
	case NEVER:
		return SUSPENDED
	default:
		return SUSPENDED
	}
}

func (c *Client) StartInstance(_ *gofr.Context, projectID, instanceName string) error {
	patchReq := &sqladmin.DatabaseInstance{
		Settings: &sqladmin.Settings{
			ActivationPolicy: ALWAYS,
		},
	}

	_, err := c.SQL.Patch(projectID, instanceName, patchReq).Do()
	if err != nil {
		return getError(err)
	}

	return nil
}

func (c *Client) StopInstance(_ *gofr.Context, projectID, instanceName string) error {
	patchReq := &sqladmin.DatabaseInstance{
		Settings: &sqladmin.Settings{
			ActivationPolicy: NEVER,
		},
	}

	_, err := c.SQL.Patch(projectID, instanceName, patchReq).Do()
	if err != nil {
		return getError(err)
	}

	return nil
}

func getError(err error) error {
	if err == nil {
		return nil
	}

	var gErr *googleapi.Error

	if errors.As(err, &gErr) {
		if gErr.Code == http.StatusConflict {
			return &ErrConflict{Message: gErr.Message}
		}
	}

	return &InternalServerError{}
}
