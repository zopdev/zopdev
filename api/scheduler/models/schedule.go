package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Schedule struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	TimeZone       string    `json:"timezone"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	ScheduleString Schedules `json:"schedule_string"`
}

type Schedules []string

func (a *Schedules) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return driver.ErrSkip
	}

	return json.Unmarshal(bytes, a)
}

func (a Schedules) Value() (driver.Value, error) {
	return json.Marshal(a)
}
