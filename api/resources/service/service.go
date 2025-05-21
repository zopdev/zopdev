package service

import (
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

type Service struct {
	gcp  GCPClient
	http HTTPClient
}

func New(gcp GCPClient, http HTTPClient) *Service {
	return &Service{gcp: gcp, http: http}
}

func (s *Service) GetResources(ctx *gofr.Context, id int64, resources []string) ([]models.Instance, error) {
	ca, err := s.http.GetCloudCredentials(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return s.getAllInstances(ctx, ca)
	}

	var instances []models.Instance

	for _, resource := range resources {
		switch resource {
		case string(SQL):
			sql, erRes := s.getAllSQLInstances(ctx, CloudDetails{
				CloudType: CloudProvider(strings.ToUpper(ca.Provider)),
				Creds:     ca.Credentials,
			})
			if erRes != nil {
				return nil, erRes
			}

			instances = append(instances, sql...)
		default:
			return nil, gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}}
		}
	}

	return instances, nil
}

func (s *Service) ChangeState(ctx *gofr.Context, resDetails ResourceDetails) error {
	ca, err := s.http.GetCloudCredentials(ctx, resDetails.CloudAccID)
	if err != nil {
		return err
	}

	switch resDetails.Type {
	case SQL:
		return s.changeSQLState(ctx, ca, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.Type"}}
	}
}
