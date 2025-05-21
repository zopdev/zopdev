package service

type (
	CloudProvider string
	ResourceType  string
)

const (
	// Cloud Providers that are currently supported in zopdev.

	GCP CloudProvider = "GCP"

	// Resource Types that are currently supported in zopdev.
	// TODO: add more resource types.

	SQL ResourceType = "sql"
	VM  ResourceType = "vm"
)

type CloudDetails struct {
	CloudType CloudProvider `json:"cloudType,omitempty"`
	Creds     any           `json:"creds"`
}

type ResourceDetails struct {
	CloudAccID int64        `json:"cloudAccID"`
	Name       string       `json:"name"`
	Type       ResourceType `json:"type"`
}
