package store

// DeploymentSpace represents a logical unit in a cloud environment.
// It contains information about a deployment space, including associated cloud account details,
// the environment ID, and timestamps for tracking the creation, update, and optional deletion of the deployment space.
type DeploymentSpace struct {
	// ID is the unique identifier for the deployment space.
	ID int64 `json:"id"`

	// CloudAccountID is the ID of the cloud account associated with the deployment space.
	CloudAccountID int64 `json:"cloudAccountId"`

	// EnvironmentID is the ID of the environment to which the deployment space belongs.
	EnvironmentID int64 `json:"environmentId"`

	// EnvironmentName is the name of the environment to which the deployment space belongs.
	EnvironmentName string `json:"environmentName"`

	// CloudAccountName is the name of the cloud account associated with the deployment space.
	CloudAccountName string `json:"cloudAccountName"`

	// Type is the type of the deployment space (e.g., production, staging).
	Type string `json:"type"`

	// CreatedAt is the timestamp of when the deployment space was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the deployment space.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the deployment space was deleted, if applicable.
	// This field is omitted if the deployment space is not deleted.
	DeletedAt string `json:"deletedAt,omitempty"`
}

// Cluster represents a cluster within a deployment space.
// A cluster is a set of computing resources within a specific region and provider.
// It includes the cluster's identifier, name, and associated metadata.
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
	// This field is omitted if the cluster is not deleted.
	DeletedAt string `json:"deletedAt,omitempty"`

	// Namespace represents the logical partition within the cluster.
	Namespace Namespace `json:"namespace"`
}

// Namespace represents a logical partition within a cluster.
// A namespace allows for isolating resources within the same cluster.
type Namespace struct {
	// Name is the name of the namespace.
	Name string `json:"Name"`
}
