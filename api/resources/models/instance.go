package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Instance struct {
	ID           int64        `json:"id"`
	Name         string       `json:"instance_name"`
	Type         string       `json:"instance_type"`
	CloudAccount CloudAccount `json:"cloud_account"`
	Region       string       `json:"region"`
	CreationTime string       `json:"creation_time"`
	Status       string       `json:"status"`
	UID          string       `json:"uid"`
	Settings     Settings     `json:"settings"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type Settings map[string]any

func (s Settings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Settings) Scan(value any) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return driver.ErrSkip
	}

	return json.Unmarshal(data, &s)
}

type CloudAccount struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}
