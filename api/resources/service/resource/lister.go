package resource

import (
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

func (s *Service) getAllInstances(ctx *gofr.Context, ca *client.CloudAccount) ([]models.Instance, error) {
	var instances []models.Instance

	// Get all SQL instances
	sql, err := s.getAllSQLInstances(ctx, CloudDetails{
		CloudType: CloudProvider(strings.ToUpper(ca.Provider)),
		Creds:     ca.Credentials,
	})
	if err != nil {
		return nil, err
	}

	instances = append(instances, sql...)

	// Get all other instances (e.g., Compute Engine, Kubernetes, etc.)
	// TODO: Implement other instance types

	return instances, nil
}

func (s *Service) getAllSQLInstances(ctx *gofr.Context, req CloudDetails) ([]models.Instance, error) {
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

	sqlClient, err := s.gcp.NewSQLClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return sqlClient.GetAllInstances(ctx, creds.ProjectID)
}
