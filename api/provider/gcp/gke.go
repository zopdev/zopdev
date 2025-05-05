// Package gcp provides an implementation of the Provider interface for interacting with
// Google Cloud Platform (GCP) resources such as GKE clusters and namespaces.
//
// It implements methods to list all clusters and namespaces in a GKE cluster using GCP credentials,
// and returns responses in the format expected by the `provider` package.
package gcp

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/zopdev/zopdev/api/provider"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"
	apiContainer "google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// GCP implements the provider.Provider interface for Google Cloud Platform.
type GCP struct {
}

// New initializes and returns a new GCP provider.
func New() provider.Provider {
	return &GCP{}
}

// ListAllClusters lists all clusters available for a given cloud account in GCP.
// It uses the GCP credentials to authenticate and fetch the cluster details.
func (g *GCP) ListAllClusters(ctx *gofr.Context, cloudAccount *provider.CloudAccount,
	credentials interface{}) (*provider.ClusterResponse, error) {
	credBody, err := g.getCredGCP(credentials)
	if err != nil {
		return nil, err
	}

	client, err := g.getClusterManagerClientGCP(ctx, credBody)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	req := &containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/-", cloudAccount.ProviderID),
	}

	resp, err := client.ListClusters(ctx, req)
	if err != nil {
		return nil, err
	}

	gkeClusters := make([]provider.Cluster, 0)

	for _, cl := range resp.Clusters {
		gkeCluster := provider.Cluster{
			Name:       cl.Name,
			Identifier: cl.Id,
			Region:     cl.Location,
			Locations:  cl.Locations,
			Type:       "deploymentSpace",
		}

		for _, nps := range cl.NodePools {
			cfg := nps.GetConfig()

			nodepool := provider.NodePool{
				MachineType: cfg.MachineType,
				NodeVersion: nps.Version,
				CurrentNode: nps.InitialNodeCount,
				NodeName:    nps.Name,
			}

			gkeCluster.NodePools = append(gkeCluster.NodePools, nodepool)
		}

		gkeClusters = append(gkeClusters, gkeCluster)
	}

	response := &provider.ClusterResponse{
		Clusters: gkeClusters,
		Next: provider.Next{
			Name: "Namespace",
			Path: fmt.Sprintf("/cloud-accounts/%v/deployment-space/namespaces", cloudAccount.ID),
			Params: map[string]string{
				"region": "region",
				"name":   "name",
			},
		},
		Metadata: provider.Metadata{
			Name: "GKE Cluster",
		},
	}

	return response, nil
}

// ListNamespace fetches namespaces from the Kubernetes API for a given GKE cluster.
func (g *GCP) ListNamespace(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAccount *provider.CloudAccount, credentials interface{}) (interface{}, error) {
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

	// Step 4: Fetch namespaces from the Kubernetes API
	apiEndpoint := fmt.Sprintf("https://%s/api/v1/namespaces", gkeCluster.Endpoint)

	namespaces, err := g.fetchNamespaces(ctx, client, credBody, apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch namespaces: %w", err)
	}

	return namespaces, nil
}

// getClusterInfo retrieves detailed information about a specific GKE cluster.
func (*GCP) getClusterInfo(ctx *gofr.Context, cluster *provider.Cluster,
	cloudAccount *provider.CloudAccount, credBody []byte) (*apiContainer.Cluster, error) {
	// Create the GCP Container service
	containerService, err := apiContainer.NewService(ctx, option.WithCredentialsJSON(credBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create container service: %w", err)
	}

	// Construct the full cluster name
	clusterFullName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		cloudAccount.ProviderID, cluster.Region, cluster.Name)

	// Get the GCP cluster details
	gkeCluster, err := containerService.Projects.Locations.Clusters.Get(clusterFullName).
		Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get GCP cluster details: %w", err)
	}

	return gkeCluster, nil
}

// createTLSConfiguredClient creates an HTTP client with custom TLS configuration using the provided CA certificate.
func (*GCP) createTLSConfiguredClient(caCertificate string) (*http.Client, error) {
	// Decode the Base64-encoded CA certificate
	caCertBytes, err := base64.StdEncoding.DecodeString(caCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CA certificate: %w", err)
	}

	// Create a CA certificate pool
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCertBytes) {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client, nil
}

// fetchNamespaces fetches Kubernetes namespaces from the specified API endpoint using the provided HTTP client.
func (*GCP) fetchNamespaces(ctx *gofr.Context, client *http.Client, credBody []byte,
	apiEndpoint string) (*provider.NamespaceResponse, error) {
	// Generate a JWT token from the credentials
	config, err := google.JWTConfigFromJSON(credBody, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %w", err)
	}

	// Create a TokenSource
	tokenSource := config.TokenSource(ctx)

	// Get a token
	token, err := tokenSource.Token()
	if err != nil {
		ctx.Errorf("failed to get token: %v", err)
		return nil, err
	}

	// Make a request to the Kubernetes API to list namespaces
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiEndpoint, http.NoBody)
	if err != nil {
		ctx.Errorf("failed to create request: %w", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle unexpected status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		ctx.Errorf("API call failed with status code %d: %s", resp.StatusCode, body)

		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var namespaceResponse struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &namespaceResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract namespace names
	namespaces := []provider.Namespace{}

	for _, item := range namespaceResponse.Items {
		namespace := provider.Namespace{
			Name: item.Metadata.Name,
			Type: "deploymentSpace.namespace",
		}

		namespaces = append(namespaces, namespace)
	}

	return &provider.NamespaceResponse{
		Options: namespaces,
		Metadata: provider.Metadata{
			Name: "namespace",
		},
	}, nil
}

// getCredGCP extracts and marshals the credentials into the appropriate format for GCP authentication.
func (*GCP) getCredGCP(credentials any) ([]byte, error) {
	var cred gcpCredentials

	credBody, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(credBody, &cred)
	if err != nil {
		return nil, err
	}

	return json.Marshal(cred)
}

// getClusterManagerClientGCP creates a client for interacting with the GKE Cluster Manager API.
func (*GCP) getClusterManagerClientGCP(ctx *gofr.Context, credentials []byte) (*container.ClusterManagerClient, error) {
	client, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, err
	}

	return client, nil
}
