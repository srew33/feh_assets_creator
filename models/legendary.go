package models

import (
	"database/sql/driver"

	"github.com/goccy/go-json"
)

type Legendary struct {
	DuoSkillID  *string `json:"duo_skill_id"`
	BonusEffect Stats   `json:"bonus_effect"`
	Kind        int64   `json:"kind"`
	Element     int64   `json:"element"`
	Bst         int64   `json:"bst"`
	PairUp      bool    `json:"pair_up"`
	AEExtra     bool    `json:"ae_extra"`
}

func (t *Legendary) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t Legendary) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}
