package resource

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

// CloudResourceProvider is the contract for cloud resource operations used by the service layer.
// Providers should implement this interface but should NOT import this package.
type CloudResourceProvider interface {
	ListResources(ctx *gofr.Context, creds any, filter models.ResourceFilter) ([]models.Resource, error)
	StartResource(ctx *gofr.Context, creds any, resource *models.Resource) error
	StopResource(ctx *gofr.Context, creds any, resource *models.Resource) error
}

type HTTPClient interface {
	GetCloudCredentials(ctx *gofr.Context, cloudAccID int64) (*client.CloudAccount, error)
	GetAllCloudAccounts(ctx *gofr.Context) ([]client.CloudAccount, error)
}

type Store interface {
	InsertResource(ctx *gofr.Context, resources *models.Resource) error
	GetResources(ctx *gofr.Context, cloudAccountID int64, resourceType []string) ([]models.Resource, error)
	UpdateStatus(ctx *gofr.Context, status string, id int64) error
	RemoveResource(ctx *gofr.Context, id int64) error
	GetResourceByID(ctx *gofr.Context, id int64) (*models.Resource, error)
}
