package gcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zopdev/zopdev/api/provider"
	"gofr.dev/pkg/gofr"

	"golang.org/x/oauth2/google"
)

func (g *GCP) ListServices(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAccount *provider.CloudAccount, credentials interface{}, namespace string) (interface{}, error) {
	// Step 1: Get GCP credentials
	credBody, err := g.getCredGCP(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Step 2: Get cluster information
	gkeCluster, err := g.getClusterInfo(ctx, cluster, cloudAccount, credBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	// Step 3: Create HTTP client with TLS configured
	client, err := g.createTLSConfiguredClient(gkeCluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS configured client: %w", err)
	}

	// Step 4: Fetch services from the Kubernetes API
	apiEndpoint := fmt.Sprintf("https://%s/api/v1/namespaces/%s/services", gkeCluster.Endpoint, namespace)

	// Parse JSON response
	var serviceResponse struct {
		Items []provider.Service `json:"items"`
	}

	err = g.fetchServices(ctx, client, credBody, apiEndpoint, &serviceResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch services: %w", err)
	}

	return &provider.ServiceResponse{
		Services: serviceResponse.Items,
		Metadata: provider.Metadata{
			Name: "services",
			Type: "kubernetes-cluster",
		},
	}, nil
}

// fetchServices fetches Kubernetes services from the specified namespace using the provided HTTP client.
func (*GCP) fetchServices(ctx *gofr.Context, client *http.Client, credBody []byte,
	apiEndpoint string, i any) error {
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

	// Make a request to the Kubernetes API to list services
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
		return errUnexpectedStatusCode{statusCode: resp.StatusCode}
	}

	if err := json.Unmarshal(body, i); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract service details

	return nil
}

func (g *GCP) GetService(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAccount *provider.CloudAccount, credentials interface{}, namespace, name string) (interface{}, error) {
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

	apiEndpoint := fmt.Sprintf("https://%s/api/v1/namespaces/%s/services/%s", gkeCluster.Endpoint, namespace, name)

	var serviceResponse provider.Service

	err = g.fetchServices(ctx, client, credBody, apiEndpoint, &serviceResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service %s: %w", name, err)
	}

	return serviceResponse, nil
}
