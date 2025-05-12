package store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

var (
	errFailedAssertion = errors.New("failed to scan JSONB: type assertion to []byte failed")
)

type Result struct {
	ID             int64       `json:"id"`
	CloudAccountID int64       `json:"cloudAccountId"`
	EvaluatedAt    time.Time   `json:"evaluatedAt"`
	RuleID         string      `json:"ruleId"`
	Result         *ResultData `json:"result"`
}

type ResultData struct {
	Data []Items `json:"items"`
}

type Items struct {
	InstanceName string `json:"instance_name"`
	Status       string `json:"status"`
	Metadata     any    `json:"metadata"`
}

func (j *ResultData) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}

	return json.Marshal(j.Data)
}

func (j *ResultData) Scan(value any) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errFailedAssertion
	}

	return json.Unmarshal(bytes, &j.Data)
}
