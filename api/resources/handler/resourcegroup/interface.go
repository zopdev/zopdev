package resourcegroup

import (
	"github.com/zopdev/zopdev/api/resources/models"
	"gofr.dev/pkg/gofr"
)

type Service interface {
	GetAllResourceGroups(ctx *gofr.Context, cloudAccID int64) ([]models.ResourceGroupData, error)
	GetResourceGroupByID(ctx *gofr.Context, cloudAccID, id int64) (*models.ResourceGroupData, error)
	CreateResourceGroup(ctx *gofr.Context, rg *models.ResourceGroup) (*models.ResourceGroupData, error)
	UpdateResourceGroup(ctx *gofr.Context, rg *models.ResourceGroup) (*models.ResourceGroupData, error)
	DeleteResourceGroup(ctx *gofr.Context, cloudAccID, id int64) error

	AddResourceToGroup(ctx *gofr.Context, groupID, resourceID int64) error
	RemoveResourceFromGroup(ctx *gofr.Context, groupID, resourceID int64) error
}
