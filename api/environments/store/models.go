package store

// Environment represents the configuration and metadata for a specific application environment.
// It includes information about the environment's name, level, associated application, and timestamps.
type Environment struct {

	// ID is the unique identifier of the environment.
	ID int64 `json:"id"`

	// ApplicationID is the unique identifier of the application to which this environment belongs.
	ApplicationID int64 `json:"applicationId"`

	// Name is the name of the environment (e.g., "development", "staging", "production").
	Name string `json:"name"`

	// CreatedAt is the timestamp indicating when the environment was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update made to the environment.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp indicating when the environment was deleted, if applicable.
	// This field is optional and may be omitted if the environment is active.
	DeletedAt string `json:"deletedAt,omitempty"`

	// Level represents the environment's hierarchical level or priority, such as an integer scale.
	Level int `json:"level"`

	// DeploymentSpace contains any additional deployment-related configuration or metadata.
	DeploymentSpace any `json:"deploymentSpace"`
}
