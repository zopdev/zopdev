package resourcegroup

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

type RGStore interface {
	GetAllResourceGroups(ctx *gofr.Context, cloudAccID int64) ([]models.ResourceGroup, error)
	GetResourceGroupByID(ctx *gofr.Context, cloudAccID, id int64) (*models.ResourceGroup, error)
	CreateResourceGroup(ctx *gofr.Context, resourceGroup *models.RGCreate) (int64, error)
	UpdateResourceGroup(ctx *gofr.Context, resourceGroup *models.RGUpdate) error
	DeleteResourceGroup(ctx *gofr.Context, id int64) error

	GetResourceIDs(ctx *gofr.Context, id int64) ([]int64, error)
	AddResourcesToGroup(ctx *gofr.Context, groupID int64, resourceID []int64) error
	RemoveResourceFromGroup(ctx *gofr.Context, groupID, resourceID int64) error
}

type ResourceService interface {
	GetByID(ctx *gofr.Context, id int64) (*models.Resource, error)
}
