package store

import "time"

type ResourceType string

type Resource struct {
	ID             int64        `json:"id"`
	UID            string       `json:"uid"`
	Name           string       `json:"name"`
	Type           ResourceType `json:"type"`
	State          string       `json:"state"`
	CloudAccountID int64        `json:"cloud_account_id"`
	CloudProvider  string       `json:"cloud_provider"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
