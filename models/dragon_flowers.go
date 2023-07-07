package models

import (
	"database/sql/driver"

	"github.com/goccy/go-json"
)

type Dragonflowers struct {
	MaxCount int64   `json:"max_count"`
	Costs    []int64 `json:"costs"`
}

func (t *Dragonflowers) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t Dragonflowers) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}
