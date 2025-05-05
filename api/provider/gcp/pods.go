package gcp

import (
	"fmt"

	"github.com/zopdev/zopdev/api/provider"

	"gofr.dev/pkg/gofr"
)

func (g *GCP) ListPods(ctx *gofr.Context, cluster *provider.Cluster,
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

	apiEndpoint := fmt.Sprintf("https://%s/api/v1/namespaces/%s/pods", gkeCluster.Endpoint, namespace)

	var podResponse struct {
		Items []provider.PodData `json:"items"`
	}

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &podResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pods: %w", err)
	}

	return &provider.Pods{
		Pods: podResponse.Items,
		Metadata: &provider.Metadata{
			Name: "pods",
			Type: "kubernetes-cluster",
		},
	}, nil
}

func (g *GCP) GetPod(ctx *gofr.Context, cluster *provider.Cluster,
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

	apiEndpoint := fmt.Sprintf("https://%s/api/v1/namespaces/%s/pods/%s", gkeCluster.Endpoint, namespace, name)

	var podResponse provider.PodData

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &podResponse)
	if err != nil {
		return nil, err
	}

	return &podResponse, nil
}
