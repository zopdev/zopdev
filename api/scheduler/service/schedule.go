package service

import (
	"errors"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/scheduler/models"
)

type Store interface {
	GetAllSchedule(ctx *gofr.Context) ([]*models.Schedule, error)
	GetScheduleByID(ctx *gofr.Context, id int64) (*models.Schedule, error)
	CreateSchedule(ctx *gofr.Context, sch *models.Schedule) (*models.Schedule, error)
	UpdateSchedule(ctx *gofr.Context, sch *models.Schedule) error
	DeleteSchedule(ctx *gofr.Context, id int64) error
}

var (
	errInternalServer = errors.New("something went wrong! Please try again later")
)

type Service struct {
	str Store
}

func New(str Store) *Service {
	return &Service{str: str}
}

func (s *Service) GetAllSchedules(ctx *gofr.Context) ([]*models.Schedule, error) {
	sch, err := s.str.GetAllSchedule(ctx)
	if err != nil {
		ctx.Logger.Errorf("error getting all schedules: %v", err)
		return nil, errInternalServer
	}

	return sch, nil
}

func (s *Service) GetSchedule(ctx *gofr.Context, id int64) (*models.Schedule, error) {
	sch, err := s.str.GetScheduleByID(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("error geting schedule: %v with id %d", err, id)
		return nil, errInternalServer
	}

	if sch == nil {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "schedule"}
	}

	return sch, nil
}

func (s *Service) CreateSchedule(ctx *gofr.Context, sch *models.Schedule) (*models.Schedule, error) {
	sch, err := s.str.CreateSchedule(ctx, sch)
	if err != nil {
		ctx.Logger.Errorf("error creating schedule: %v", err)
		return nil, errInternalServer
	}

	return sch, nil
}

func (s *Service) UpdateSchedule(ctx *gofr.Context, sch *models.Schedule) error {
	err := s.str.UpdateSchedule(ctx, sch)
	if err != nil {
		ctx.Logger.Errorf("error updating schedule: %v", err)
		return errInternalServer
	}

	return nil
}

func (s *Service) DeleteSchedule(ctx *gofr.Context, id int64) error {
	if id == 0 {
		return gofrHttp.ErrorMissingParam{Params: []string{"id"}}
	}

	err := s.str.DeleteSchedule(ctx, id)
	if err != nil {
		ctx.Logger.Errorf("error deleting schedule: %v", err)
		return errInternalServer
	}

	return nil
}
