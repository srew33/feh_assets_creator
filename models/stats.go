package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type Stat struct {
	Key   int
	Value int64
}

type Stats struct {
	HP  int64 `json:"hp"`
	Atk int64 `json:"atk"`
	Spd int64 `json:"spd"`
	Def int64 `json:"def"`
	Res int64 `json:"res"`
}

func (b *Stats) ToList() []Stat {
	return []Stat{
		{Key: HP, Value: b.HP},
		{Key: Atk, Value: b.Atk},
		{Key: Spd, Value: b.Spd},
		{Key: Def, Value: b.Def},
		{Key: Res, Value: b.Res},
	}
}
func (b *Stats) FromList(stats []Stat) {
	b.HP = stats[0].Value
	b.Atk = stats[1].Value
	b.Spd = stats[2].Value
	b.Def = stats[3].Value
	b.Res = stats[4].Value
}

func (b *Stats) Sum() int64 {
	return b.HP + b.Atk + b.Spd + b.Def + b.Res
}
func (b *Stats) IsDefault() bool {
	return b.HP == 0 && b.Atk == 0 && b.Spd == 0 && b.Def == 0 && b.Res == 0
}
func (b *Stats) String() string {
	m := map[string]int64{
		"hp":  b.HP,
		"atk": b.Atk,
		"spd": b.Spd,
		"def": b.Def,
		"res": b.Res,
	}
	by := ""
	for k, v := range m {
		if v != 0 {
			// by.WriteString(fmt.Sprintf("%s+%d ", strings.ToUpper(k), v))
			by += fmt.Sprintf("%s+%d ", strings.ToUpper(k), v)
		}
	}
	return by
}

const (
	Null = iota
	HP
	Atk
	Spd
	Def
	Res
)

var DisadvantageRates = []int{10, 30, 50, 75, 95}

func (t *Stats) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t Stats) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}
