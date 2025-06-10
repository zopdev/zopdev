package resource

import (
	"strings"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

func (s *Service) getAllInstances(ctx *gofr.Context, ca *client.CloudAccount) ([]models.Resource, error) {
	var instances []models.Resource

	type result struct {
		instances []models.Resource
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

func (s *Service) getAllSQLInstances(ctx *gofr.Context, req CloudDetails) ([]models.Resource, error) {
	switch req.CloudType {
	case GCP:
		return s.getGCPSQLInstances(ctx, req.Creds)
	case AWS:
		return s.getAWSRDSInstances(ctx, req.Creds)
	default:
		// We are not returning any error because the sync process is completely internal, works on the cloud Account ID,
		// if we are getting an unknown cloud type, then this feature is not implemented and we simply return nil.
		return nil, nil
	}
}

func (s *Service) getGCPSQLInstances(ctx *gofr.Context, creds any) ([]models.Resource, error) {
	filter := models.ResourceFilter{
		ResourceTypes: []string{"SQL"},
	}

	return s.gcp.ListResources(ctx, creds, filter)
}

func (s *Service) getAWSRDSInstances(ctx *gofr.Context, creds any) ([]models.Resource, error) {
	filter := models.ResourceFilter{
		ResourceTypes: []string{"RDS"},
	}

	return s.aws.ListResources(ctx, creds, filter)
}
