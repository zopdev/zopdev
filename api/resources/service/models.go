package service

type CloudProvider int32

const (
	GCP   CloudProvider = 0
	AWS   CloudProvider = 1
	AZURE CloudProvider = 2
)

type Request struct {
	CloudType CloudProvider `json:"cloudType,omitempty"`
	Creds     any           `json:"creds"`
}
