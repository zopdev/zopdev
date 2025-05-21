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

	ALL ResourceType = "all"
	SQL ResourceType = "sql"
	VM  ResourceType = "vm"
)

type Request struct {
	CloudType CloudProvider `json:"cloudType,omitempty"`
	Creds     any           `json:"creds"`
}
