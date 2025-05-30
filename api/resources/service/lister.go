package service

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

	type result struct {
		instances []models.Instance
		err       error
	}

	sqlCh := make(chan result, 1)
	computeCh := make(chan result, 1)

	// Fetch SQL instances concurrently
	go func() {
		sql, err := s.getAllSQLInstances(ctx, CloudDetails{
			CloudType: CloudProvider(strings.ToUpper(ca.Provider)),
			Creds:     ca.Credentials,
		})

		for i := range sql {
			sql[i].CloudAccount.ID = ca.ID
			sql[i].CloudAccount.Type = ca.Provider
		}

		sqlCh <- result{sql, err}
	}()

	// Fetch compute instances concurrently
	go func() {
		computeInstances, err := s.getALLComputeInstances(ctx, CloudDetails{
			CloudType: CloudProvider(strings.ToUpper(ca.Provider)),
			Creds:     ca.Credentials,
		})

		for i := range computeInstances {
			computeInstances[i].CloudAccount.ID = ca.ID
			computeInstances[i].CloudAccount.Type = ca.Provider
		}

		computeCh <- result{computeInstances, err}
	}()

	var sqlRes, computeRes result
	for i := 0; i < 2; i++ {
		select {
		case r := <-sqlCh:
			sqlRes = r
		case r := <-computeCh:
			computeRes = r
		}
	}

	if sqlRes.err != nil {
		return nil, sqlRes.err
	}
	if computeRes.err != nil {
		return nil, computeRes.err
	}

	instances = append(instances, sqlRes.instances...)
	instances = append(instances, computeRes.instances...)

	// Get all other instances (e.g., Compute Engine, Kubernetes, etc.)
	// TODO: Implement other instance types

	return instances, nil
}

func (s *Service) getAllSQLInstances(ctx *gofr.Context, req CloudDetails) ([]models.Instance, error) {
	switch req.CloudType {
	case GCP:
		return s.getGCPSQLInstances(ctx, req.Creds)
	case AWS:
		return s.getAWSRDSInstances(ctx, req.Creds)
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

func (s *Service) getAWSRDSInstances(ctx *gofr.Context, cred any) ([]models.Instance, error) {
	awsRDSClient, err := s.aws.NewRDSClient(ctx, cred)
	if err != nil {
		return nil, err
	}

	return awsRDSClient.GetAllInstances(ctx)
}
