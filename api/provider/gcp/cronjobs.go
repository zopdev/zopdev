package gcp

import (
	"fmt"

	"github.com/zopdev/zopdev/api/provider"

	"gofr.dev/pkg/gofr"
)

func (g *GCP) ListCronJobs(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAccount *provider.CloudAccount, credentials interface{}, namespace string) (interface{}, error) {
	credBody, err := g.getCredGCP(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	gkeCluster, err := g.getClusterInfo(ctx, cluster, cloudAccount, credBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	client, err := g.createTLSConfiguredClient(gkeCluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS configured client: %w", err)
	}

	apiEndpoint := fmt.Sprintf("https://%s/apis/batch/v1/namespaces/%s/cronjobs", gkeCluster.Endpoint, namespace)

	var cronJobResponse struct {
		Items []provider.CronJobData `json:"items"`
	}

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &cronJobResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cronjobs: %w", err)
	}

	return &provider.CronJobs{
		CronJobs: cronJobResponse.Items,
		Metadata: &provider.Metadata{
			Name: "cronjobs",
			Type: "kubernetes-cluster",
		},
	}, nil
}

func (g *GCP) GetCronJob(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAcc *provider.CloudAccount, creds any, namespace, name string) (any, error) {
	credBody, err := g.getCredGCP(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	gkeCluster, err := g.getClusterInfo(ctx, cluster, cloudAcc, credBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	client, err := g.createTLSConfiguredClient(gkeCluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS configured client: %w", err)
	}

	apiEndpoint := fmt.Sprintf("https://%s/apis/batch/v1/namespaces/%s/cronjobs/%s", gkeCluster.Endpoint, namespace, name)

	var cronJobResponse provider.CronJobData

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &cronJobResponse)
	if err != nil {
		return nil, err
	}

	return &cronJobResponse, nil
}
