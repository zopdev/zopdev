package resource

import (
	"strings"

	"github.com/zopdev/zopdev/api/resources/models"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

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

func (s *Service) changeGCPSQL(ctx *gofr.Context, cred any, resDetails ResourceDetails) error {
	creds, err := s.gcp.NewGoogleCredentials(ctx, cred, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return err
	}

	sqlClient, err := s.gcp.NewSQLClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return err
	}

	switch resDetails.State {
	case START:
		return sqlClient.StartInstance(ctx, creds.ProjectID, resDetails.Name)
	case SUSPEND:
		return sqlClient.StopInstance(ctx, creds.ProjectID, resDetails.Name)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.State"}}
	}
}

func (s *Service) changeAWSRDS(ctx *gofr.Context, cred any, state ResourceState, resDetails *models.Resource) error {
	cl, err := s.aws.NewRDSClient(ctx, cred)
	if err != nil {
		return err
	}

	switch state {
	case START:
		return cl.StartInstance(ctx, resDetails)
	case SUSPEND:
		return cl.StopInstance(ctx, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.state"}}
	}
}
