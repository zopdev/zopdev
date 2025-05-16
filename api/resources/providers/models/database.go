package models

type SQLInstance struct {
	Name         string `json:"instance_name"`
	ProjectID    string `json:"project_id"`
	Region       string `json:"region"`
	Zone         string `json:"zone"`
	Version      string `json:"version"`
	CreationTime string `json:"creation_time"`
}
