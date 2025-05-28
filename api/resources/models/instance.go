package models

import "time"

type Instance struct {
	ID           int64          `json:"id"`
	Name         string         `json:"instance_name"`
	Type         string         `json:"instance_type"`
	CloudAccount CloudAccount   `json:"cloud_account"`
	Region       string         `json:"region"`
	CreationTime string         `json:"creation_time"`
	Status       string         `json:"status"`
	UID          string         `json:"uid"`
	Settings     map[string]any `json:"settings"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CloudAccount struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}
