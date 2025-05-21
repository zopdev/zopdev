package service

import (
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

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
	switch resDetails.State {
	case START:
		return s.start(ctx, resDetails)
	case SUSPEND:
		return s.stop(ctx, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"state"}}
	}
}
