package gcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"

	"golang.org/x/oauth2/google"

	"github.com/zopdev/zopdev/api/provider"
)

func (g *GCP) ListDeployments(ctx *gofr.Context, cluster *provider.Cluster,
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

	apiEndpoint := fmt.Sprintf("https://%s/apis/apps/v1/namespaces/%s/deployments", gkeCluster.Endpoint, namespace)

	// Parse JSON response
	var depResponse struct {
		Items []provider.DeploymentData `json:"items"`
	}

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &depResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments: %w", err)
	}

	return &provider.Deployments{
		Deployments: depResponse.Items,
		Metadata: &provider.Metadata{
			Name: "deployments",
			Type: "kubernetes-cluster",
		},
	}, nil
}

// fetchDeployments fetches Kubernetes deployments from the specified namespace using the provided HTTP client.
func (*GCP) fetchDeployments(ctx *gofr.Context, client *http.Client, credBody []byte,
	apiEndpoint string, depREsp any) error {
	// Generate a JWT token from the credentials
	config, err := google.JWTConfigFromJSON(credBody, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return fmt.Errorf("failed to create JWT config: %w", err)
	}

	// Create a TokenSource
	tokenSource := config.TokenSource(ctx)

	// Get a token
	token, err := tokenSource.Token()
	if err != nil {
		ctx.Errorf("failed to get token: %v", err)
		return err
	}

	// Make a request to the Kubernetes API to list deployments
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiEndpoint, http.NoBody)
	if err != nil {
		ctx.Errorf("failed to create request: %w", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle unexpected status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		ctx.Errorf("API call failed with status code %d: %s", resp.StatusCode, body)

		return errUnexpectedStatusCode{statusCode: resp.StatusCode}
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, depREsp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

func (g *GCP) GetDeployment(ctx *gofr.Context, cluster *provider.Cluster,
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

	apiEndpoint := fmt.Sprintf("https://%s/apis/apps/v1/namespaces/%s/deployments/%s", gkeCluster.Endpoint, namespace, name)

	var depResponse provider.DeploymentData

	err = g.fetchDeployments(ctx, client, credBody, apiEndpoint, &depResponse)
	if err != nil {
		return nil, err
	}

	return &depResponse, nil
}
