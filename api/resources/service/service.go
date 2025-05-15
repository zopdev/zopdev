package service

import (
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Service struct {
	gcp GCPClient
}

func New() *Service { return &Service{} }

func (s *Service) GetAllSQLInstances(ctx *gofr.Context, req Request) ([]models.SQLInstance, error) {
	switch req.CloudType {
	case GCP:
		return s.getGCPSQLInstances(ctx, req.Creds)
	default:
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}}
	}
}

func (s *Service) getGCPSQLInstances(ctx *gofr.Context, cred any) ([]models.SQLInstance, error) {
	creds, err := s.gcp.NewGoogleCredentials(ctx, cred, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	sqlClient, err := s.gcp.NewSQLInstanceLister(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return sqlClient.GetAllInstances(ctx, creds.ProjectID)
}
