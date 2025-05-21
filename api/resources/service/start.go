package service

import (
	"strings"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/audit/client"
)

// Start - starts the resource based on the resource type and cloud account.
func (s *Service) Start(ctx *gofr.Context, resDetails ResourceDetails) error {
	ca, err := client.GetCloudCredentials(ctx, resDetails.CloudAccID)
	if err != nil {
		return err
	}

	switch resDetails.Type {
	case SQL:
		return s.startSQL(ctx, ca, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"type"}}
	}
}

func (s *Service) startSQL(ctx *gofr.Context, ca *client.CloudAccount, resDetails ResourceDetails) error {
	var err error

	switch strings.ToUpper(ca.Provider) {
	case string(GCP):
		err = s.startGCPSQL(ctx, ca.Credentials, resDetails)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) startGCPSQL(ctx *gofr.Context, cred any, resDetails ResourceDetails) error {
	creds, err := s.gcp.NewGoogleCredentials(ctx, cred, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return err
	}

	sqlClient, err := s.gcp.NewSQLClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return err
	}

	return sqlClient.StartInstance(ctx, creds.ProjectID, resDetails.Name)
}
