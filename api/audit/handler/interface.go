package handler

import (
	"github.com/zopdev/zopdev/api/audit/store"
	"gofr.dev/pkg/gofr"
)

type Service interface {
	RunById(ctx *gofr.Context, ruleID string, cloudAccId int64) (*store.Result, error)
	RunByCategory(ctx *gofr.Context, category string, cloudAccId int64) ([]*store.Result, error)
	RunAll(ctx *gofr.Context, cloudAccId int64) (map[string][]*store.Result, error)

	GetResultById(ctx *gofr.Context, cloudAccId int64, ruleId string) (*store.Result, error)
	GetResultByCategory(ctx *gofr.Context, cloudAccId int64) (map[string][]*store.Result, error)
	GetResultByAll(ctx *gofr.Context, cloudAccId int64) ([]*store.Result, error)
}
