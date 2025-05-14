package sql

import (
	"context"

	"github.com/zopdev/zopdev/api/resources/models"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type Instance struct {
	client SQLAdmin
}

func NewInstance(client SQLAdmin) *Instance {
	return &Instance{
		client: client,
	}
}

func getCaller(project string) Caller {
	admin, err := sqladmin.NewService(context.Background())
	if err != nil {
		return nil
	}

	itr := admin.Instances.List(project)
	if err != nil {
		return nil
	}

	return itr
}

func GetAllInstance(getCaller func(projectId string) Caller, projectId string) ([]models.SQLInstance, error) {
	itr, err := getCaller(projectId).Do()
	if err != nil {
		return nil, err
	}

	var instances []models.SQLInstance
	for _, instance := range itr.Items {
		instances = append(instances, models.SQLInstance{
			InstanceName: instance.Name,
			Region:       instance.Region,
			InstanceType: instance.DatabaseVersion,
		})
	}

	return instances, nil
}
