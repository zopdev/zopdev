package service

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/audit/client"
	"github.com/zopdev/zopdev/api/audit/store"
)

// Rule is an interface that defines the methods that a rule must implement.
// It is used to execute the rule and get the result of the rule execution.
// Rule is cloud agnostic and can be used for any cloud provider, the rule should implement the logic based on the cloud provider.
type Rule interface {
	GetCategory() string
	GetName() string
	Execute(ctx *gofr.Context, ca *client.CloudAccount) ([]store.Items, error)
}

type Store interface {
	GetLastRun(ctx *gofr.Context, cloudAccID int64, rule string) (*store.Result, error)
	UpdateResult(ctx *gofr.Context, result *store.Result) error
	CreatePending(ctx *gofr.Context, result *store.Result) (*store.Result, error)
}
