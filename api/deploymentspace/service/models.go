package service

import (
	"github.com/zopdev/zopdev/api/deploymentspace/store"
	"github.com/zopdev/zopdev/api/provider"
)

// DeploymentSpaceResp represents the response structure for a deployment space.
// It contains the details of the deployment space and the associated cluster.
type DeploymentSpaceResp struct {
	// DeploymentSpace is the deployment space object containing details like CloudAccount, Type, and DeploymentSpace.
	DeploymentSpace *store.DeploymentSpace `json:"deploymentSpace"`

	// Cluster is the associated cluster within the deployment space.
	Cluster *store.Cluster `json:"cluster"`
}

// DeploymentSpace represents a deployment space in the service layer, including details about the cloud account,
// type, and deployment space itself.
type DeploymentSpace struct {
	// CloudAccount represents the cloud account associated with the deployment space.
	CloudAccount CloudAccount `json:"cloudAccount"`

	// Type represents the type of the deployment space (e.g., production, staging).
	Type Type `json:"type"`

	// DeploymentSpace represents the deployment space object itself.
	// It may include the actual deployment space model or relevant data about the space.
	DeploymentSpace interface{} `json:"deploymentSpace"`
}

// Type represents the type of the deployment space, such as the environment or the specific function the space serves.
type Type struct {
	// Name is the name of the deployment space type (e.g., production, staging).
	Name string `json:"name"`
}

// CloudAccount represents a cloud account with necessary attributes.
type CloudAccount struct {
	// ID is a unique identifier for the cloud account.
	ID int64 `json:"id,omitempty"`

	// Name is the name of the cloud account.
	Name string `json:"name"`

	// Provider is the name of the cloud service provider.
	Provider string `json:"provider"`

	// ProviderID is the identifier for the provider account.
	ProviderID string `json:"providerId"`

	// CreatedAt is the timestamp of when the cloud account was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the cloud account.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the cloud account was deleted, if applicable.
	DeletedAt string `json:"deletedAt,omitempty"`

	// ProviderDetails contains additional details specific to the provider.
	ProviderDetails interface{} `json:"providerDetails"`

	// Credentials hold authentication information for access to the provider.
	Credentials interface{} `json:"credentials,omitempty"`
}

type clusterConfigs struct {
	deploymentSpace *store.DeploymentSpace
	cluster         *provider.Cluster
	cloudAccount    *provider.CloudAccount
	credentials     interface{}
	Namespace       string
}
