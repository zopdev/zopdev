package resourcegroup

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

type Service interface {
	GetAllResourceGroups(ctx *gofr.Context, cloudAccID int64) ([]models.ResourceGroupData, error)
	GetResourceGroupByID(ctx *gofr.Context, cloudAccID, id int64) (*models.ResourceGroupData, error)
	CreateResourceGroup(ctx *gofr.Context, rg *models.RGCreate) (*models.ResourceGroupData, error)
	UpdateResourceGroup(ctx *gofr.Context, rg *models.RGUpdate) (*models.ResourceGroupData, error)
	DeleteResourceGroup(ctx *gofr.Context, cloudAccID, id int64) error
}
