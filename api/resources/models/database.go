package models

type SQLInstance struct {
	InstanceName string `json:"instance_name"`
	ProjectID    string `json:"project_id"`
	Region       string `json:"region"`
	Zone         string `json:"zone"`
	InstanceType string `json:"instance_type"`
}
