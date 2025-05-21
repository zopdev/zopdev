package service

import (
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
)

func (s *Service) changeSQLState(ctx *gofr.Context, ca *client.CloudAccount, resDetails ResourceDetails) error {
	var err error

	switch strings.ToUpper(ca.Provider) {
	case string(GCP):
		err = s.changeGCPSQL(ctx, ca.Credentials, resDetails)
	}

	if err != nil {
		return err
	}

	return nil
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
