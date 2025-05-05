package service

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/zopdev/zopdev/api/environments/service"

	"github.com/zopdev/zopdev/api/applications/store"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type Service struct {
	store              store.ApplicationStore
	environmentService service.EnvironmentService
}

func New(str store.ApplicationStore, envSvc service.EnvironmentService) ApplicationService {
	return &Service{store: str, environmentService: envSvc}
}

func (s *Service) AddApplication(ctx *gofr.Context, application *store.Application) (*store.Application, error) {
	tempApplication, err := s.store.GetApplicationByName(ctx, application.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if tempApplication != nil {
		return nil, http.ErrorEntityAlreadyExist{}
	}

	environments := application.Environments

	application, err = s.store.InsertApplication(ctx, application)
	if err != nil {
		return nil, err
	}

	if len(environments) == 0 {
		environments = append(environments, store.Environment{
			Name:  "default",
			Level: 1,
		})
	}

	for i := range environments {
		environments[i].ApplicationID = application.ID

		environment, err := s.store.InsertEnvironment(ctx, &environments[i])
		if err != nil {
			return nil, err
		}

		environments[i] = *environment
	}

	return application, nil
}

func (s *Service) FetchAllApplications(ctx *gofr.Context) ([]store.Application, error) {
	applications, err := s.store.GetALLApplications(ctx)
	if err != nil {
		return nil, err
	}

	for i := range applications {
		environments, err := s.environmentService.FetchAll(ctx, int(applications[i].ID))
		if err != nil {
			return nil, err
		}

		bytes, err := json.Marshal(environments)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bytes, &applications[i].Environments)
		if err != nil {
			return nil, err
		}
	}

	return applications, nil
}

func (s *Service) GetApplication(ctx *gofr.Context, id int) (*store.Application, error) {
	application, err := s.store.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	environments, err := s.environmentService.FetchAll(ctx, id)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(environments)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &application.Environments)
	if err != nil {
		return nil, err
	}

	return application, nil
}
