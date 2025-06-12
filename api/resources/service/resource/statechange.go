package resource

import (
	"strings"

	"github.com/zopdev/zopdev/api/resources/models"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/client"
)

func (s *Service) changeSQLState(ctx *gofr.Context, ca *client.CloudAccount, resDetails ResourceDetails) error {
	var err error

	provider := strings.ToUpper(ca.Provider)

	switch provider {
	case string(GCP):
		err = s.changeGCPSQL(ctx, ca.Credentials, resDetails)
		if err != nil {
			return err
		}

		return nil
	case string(AWS):
		resource, err := s.store.GetResourceByID(ctx, resDetails.ID)
		if err != nil {
			return err
		}

		err = s.changeAWSRDS(ctx, ca.Credentials, resDetails.State, resource)
		if err != nil {
			return err
		}

		return nil
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"cloud provider"}}
	}
}

func (s *Service) changeGCPSQL(ctx *gofr.Context, creds any, resDetails ResourceDetails) error {
	resource := models.Resource{
		ID:   resDetails.ID,
		Name: resDetails.Name,
		Type: string(resDetails.Type),
		// Add other fields as needed
	}

	switch resDetails.State {
	case START:
		return s.gcp.StartResource(ctx, creds, &resource)
	case SUSPEND:
		return s.gcp.StopResource(ctx, creds, &resource)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.State"}}
	}
}

func (s *Service) changeAWSRDS(ctx *gofr.Context, creds any, state ResourceState, resDetails *models.Resource) error {
	switch state {
	case START:
		return s.aws.StartResource(ctx, creds, resDetails)
	case SUSPEND:
		return s.aws.StopResource(ctx, creds, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.state"}}
	}
}
