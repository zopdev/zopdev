package store

import "gofr.dev/pkg/gofr"

type CloudAccountStore interface {
	InsertCloudAccount(ctx *gofr.Context, config *CloudAccount) (*CloudAccount, error)
	GetALLCloudAccounts(ctx *gofr.Context) ([]CloudAccount, error)
	GetCloudAccountByProvider(ctx *gofr.Context, providerType, providerID string) (*CloudAccount, error)
	GetCloudAccountByID(ctx *gofr.Context, cloudAccountID int64) (*CloudAccount, error)
	GetCredentials(ctx *gofr.Context, cloudAccountID int64) (interface{}, error)
}
