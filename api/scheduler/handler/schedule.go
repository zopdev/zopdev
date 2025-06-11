package handler

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/scheduler/models"
)

type Service interface {
	GetAllSchedules(ctx *gofr.Context) ([]*models.Schedule, error)
	GetSchedule(ctx *gofr.Context, resID int64) (*models.Schedule, error)
	CreateSchedule(ctx *gofr.Context, sch *models.Schedule) (*models.Schedule, error)
	UpdateSchedule(ctx *gofr.Context, sch *models.Schedule) error
	DeleteSchedule(ctx *gofr.Context, id int64) error
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAllSchedules(ctx *gofr.Context) (any, error) {
	schedules, err := h.svc.GetAllSchedules(ctx)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (h *Handler) GetSchedule(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	resID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	schedule, err := h.svc.GetSchedule(ctx, resID)
	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "schedule"}
	}

	return schedule, nil
}

func (h *Handler) CreateSchedule(ctx *gofr.Context) (any, error) {
	var sch models.Schedule

	if err := ctx.Bind(&sch); err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	newSchedule, err := h.svc.CreateSchedule(ctx, &sch)
	if err != nil {
		return nil, err
	}

	return newSchedule, nil
}

func (h *Handler) UpdateSchedule(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	scheduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	var sch models.Schedule

	err = ctx.Bind(&sch)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	sch.ID = scheduleID

	if err := h.svc.UpdateSchedule(ctx, &sch); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) DeleteSchedule(ctx *gofr.Context) (any, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return nil, gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	scheduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	err = h.svc.DeleteSchedule(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
