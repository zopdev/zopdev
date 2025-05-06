package service

import (
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/audit/store"
)

type Store interface {
	GetLastRun(ctx *gofr.Context, cloudAccID int64, rule string) (*store.Result, error)
	UpdateResult(ctx *gofr.Context, result *store.Result) error
	CreatePending(ctx *gofr.Context, result *store.Result) (*store.Result, error)
}
