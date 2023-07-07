package models

import (
	"database/sql/driver"

	"github.com/goccy/go-json"
)

type TransElement struct {
	Key string  `json:"key"`
	Val *string `json:"value"`
}

func (t *TransElement) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t TransElement) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}
