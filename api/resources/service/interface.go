package service

import (
	"context"

	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/providers/gcp"
)

type GCPClient interface {
	NewGoogleCredentials(ctx context.Context, cred any, scopes ...string) (*google.Credentials, error)
	NewSQLClient(ctx context.Context, opts ...option.ClientOption) (gcp.SQLClient, error)
}

type HTTPClient interface {
	GetCloudCredentials(ctx *gofr.Context, cloudAccID int64) (*client.CloudAccount, error)
	GetAllCloudAccounts(ctx *gofr.Context) ([]client.CloudAccount, error)
}

type Store interface {
	InsertResource(ctx *gofr.Context, resources *models.Instance) error
	GetResources(ctx *gofr.Context, cloudAccountID int64, resourceType []string) ([]models.Instance, error)
	UpdateStatus(ctx *gofr.Context, status string, id int64) error
	RemoveResource(ctx *gofr.Context, id int64) error
	GetResourceByID(ctx *gofr.Context, id int64) (*models.Instance, error)
}
