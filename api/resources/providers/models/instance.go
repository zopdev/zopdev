package models

type Instance struct {
	Name         string         `json:"instance_name"`
	Type         string         `json:"instance_type"`
	ProviderID   string         `json:"provider_id"`
	Region       string         `json:"region"`
	CreationTime string         `json:"creation_time"`
	Status       string         `json:"status"`
	UID          string         `json:"uid"`
	Settings     map[string]any `json:"settings"`
}

type Provider struct {
	ID int64 `json:"id"`
}
