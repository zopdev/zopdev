/*
Package store provides types for managing clusters and their associated namespaces.
It defines the structure for a `Cluster`, which represents a logical unit of computing resources,
and the `Namespace` within the cluster.
*/
package store

// Cluster represents a computing cluster with a unique identifier, deployment space, and configuration details.
// It includes information about the cluster's name, region, provider, and timestamps for tracking creation, update,
// and optional deletion. Each cluster is associated with a `Namespace`.
type Cluster struct {
	// DeploymentSpaceID is the unique identifier of the deployment space to which the cluster belongs.
	DeploymentSpaceID int64 `json:"deploymentSpaceId"`

	// ID is the unique identifier of the cluster.
	ID int64 `json:"id"`

	// Identifier is a unique identifier for the cluster, typically provided by the cloud provider.
	Identifier string `json:"identifier"`

	// Name is the name of the cluster.
	Name string `json:"name"`

	// Region is the geographical region where the cluster is deployed.
	Region string `json:"region"`

	// Provider is the cloud provider hosting the cluster (e.g., AWS, GCP, Azure).
	Provider string `json:"provider"`

	// ProviderID is the unique identifier of the cluster from the cloud provider's perspective.
	ProviderID string `json:"providerId"`

	// CreatedAt is the timestamp of when the cluster was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the cluster.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the cluster was deleted, if applicable.
	DeletedAt string `json:"deletedAt,omitempty"`

	// Namespace represents the namespace associated with the cluster.
	Namespace Namespace `json:"namespace"`
}

// Namespace represents a logical partition within a cluster.
// It typically defines an isolated environment for resources within the cluster.
type Namespace struct {
	// Name is the name of the namespace.
	Name string `json:"name"`
}
