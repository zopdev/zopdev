package service

import (
	"github.com/pkg/errors"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Service struct {
	gcp GCPClient
}

func New(gcp GCPClient) *Service {
	return &Service{gcp: gcp}
}

func (s *Service) GetResources(ctx *gofr.Context, id int64, resources []string) ([]models.Instance, error) {
	ca, err := client.GetCloudCredentials(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return s.getAllInstances(ctx, ca.Credentials)
	}

	var (
		instances []models.Instance
		er        error
	)

	for _, resource := range resources {
		switch resource {
		case string(SQL):
			sql, erRes := s.GetAllSQLInstances(ctx, Request{
				CloudType: CloudProvider(ca.Provider),
				Creds:     ca.Credentials,
			})
			if erRes != nil {
				er = errors.Wrap(er, erRes.Error())
			}

			instances = append(instances, sql...)
		default:
			// Return all computed instances till now
			er = errors.Wrap(er, "Unsupported resource type: "+resource)
		}
	}

	return instances, er
}

func (s *Service) getAllInstances(ctx *gofr.Context, cred any) ([]models.Instance, error) {
	var instances []models.Instance

	// Get all SQL instances
	sql, err := s.getGCPSQLInstances(ctx, cred)
	if err != nil {
		return nil, err
	}

	instances = append(instances, sql...)

	// Get all other instances (e.g., Compute Engine, Kubernetes, etc.)
	// TODO: Implement other instance types

	return instances, nil
}

func (s *Service) GetAllSQLInstances(ctx *gofr.Context, req Request) ([]models.Instance, error) {
	switch req.CloudType {
	case GCP:
		return s.getGCPSQLInstances(ctx, req.Creds)
	default:
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}}
	}
}

func (s *Service) getGCPSQLInstances(ctx *gofr.Context, cred any) ([]models.Instance, error) {
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
