package handler

import (
	"github.com/zopdev/zopdev/api/audit/store"
	"gofr.dev/pkg/gofr"
)

type Service interface {
	RunByID(ctx *gofr.Context, ruleID string, cloudAccID int64) (*store.Result, error)
	RunByCategory(ctx *gofr.Context, category string, cloudAccID int64) ([]*store.Result, error)
	RunAll(ctx *gofr.Context, cloudAccID int64) (map[string][]*store.Result, error)

	GetResultByID(ctx *gofr.Context, cloudAccID int64, ruleID string) (*store.Result, error)
	GetAllResults(ctx *gofr.Context, cloudAccID int64) (map[string][]*store.Result, error)
}
