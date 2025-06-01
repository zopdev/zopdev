package service

type (
	CloudProvider string
	ResourceType  string

	ResourceState string
)

const (
	// Cloud Providers that are currently supported in zopdev.

	GCP CloudProvider = "GCP"
	AWS CloudProvider = "AWS"

	// Resource Types that are currently supported in zopdev.
	// TODO: add more resource types.

	SQL ResourceType = "SQL"

	// Resource State constants.

	START   ResourceState = "START"
	SUSPEND ResourceState = "SUSPEND"

	RUNNING = "RUNNING"
	STOPPED = "STOPPED"
)

type CloudDetails struct {
	CloudType CloudProvider `json:"cloudType,omitempty"`
	Creds     any           `json:"creds"`
}

type ResourceDetails struct {
	ID         int64         `json:"id"`
	CloudAccID int64         `json:"cloudAccID"`
	Name       string        `json:"name"`
	Type       ResourceType  `json:"type"`
	State      ResourceState `json:"state"`
}
