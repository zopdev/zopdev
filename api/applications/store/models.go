/*
Package store provides types for managing applications and their environments.
It defines the structure for an `Application` and its associated `Environment`.

Applications represent a logical grouping of resources, while Environments
represent specific configurations or stages (e.g., development, staging, production)
within an Application.
*/
package store

// Application represents an application with a unique ID, name, and associated environments.
// It also includes timestamps for tracking the creation, update, and optional deletion of the application.
type Application struct {
	// ID is the unique identifier of the application.
	ID int64 `json:"id"`

	// Name is the name of the application.
	Name string `json:"name"`

	// Environments is a list of environments associated with the application.
	Environments []Environment `json:"environments"`

	// CreatedAt is the timestamp of when the application was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the application.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the application was deleted, if applicable.
	DeletedAt string `json:"deletedAt,omitempty"`
}

// Environment represents a specific environment within an application.
// Environments have unique identifiers, a name, a hierarchical level, and a reference to their parent application.
// Each environment also contains deployment-specific details.
type Environment struct {
	// ID is the unique identifier of the environment.
	ID int64 `json:"id"`

	// Name is the name of the environment.
	Name string `json:"name"`

	// Level represents the hierarchical level of the environment.
	// For example, 1 for development, 2 for staging, and 3 for production.
	Level int `json:"level"`

	// ApplicationID is the ID of the application to which the environment belongs.
	ApplicationID int64 `json:"applicationID"`

	// DeploymentSpace contains configuration details specific to the deployment in this environment.
	// The type `any` allows for flexibility in defining deployment-specific data.
	DeploymentSpace any `json:"deploymentSpace"`

	// CreatedAt is the timestamp of when the environment was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the environment.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the environment was deleted, if applicable.
	DeletedAt string `json:"deletedAt,omitempty"`
}
