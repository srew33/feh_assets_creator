package models

import (
	"database/sql/driver"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goccy/go-json"
)

type C comparable

// 直接使用comparable，类型断言对value不起作用，编译器会报错
type Set[T C] struct {
	mapset.Set[T]
}

func (t *Set[T]) Scan(value interface{}) error {
	bytesValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("")
	}

	t.Set = mapset.NewSet[T]()

	return json.Unmarshal([]byte(bytesValue), t.Set)
}

func (t Set[T]) Value() (driver.Value, error) {
	if t.Set == nil {
		t.Set = mapset.NewSet[T]()
	}
	by, err := json.Marshal(t.Set)
	// fmt.Println(by)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}

func (t *Set[string]) MarshalJSON() ([]byte, error) {

	b, err := t.Set.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Set[string]) UnmarshalJSON(data []byte) error {

	var s = mapset.NewSet[string]()

	err := s.UnmarshalJSON(data)
	t.Set = s

	return err
}

type SliceI64 []int64

func (t *SliceI64) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t SliceI64) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}

type SliceStr []*string

func (t *SliceStr) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t SliceStr) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}

type StringBool bool

func (t *StringBool) Scan(value interface{}) error {
	v, _ := value.(string)
	var d string
	err := json.Unmarshal([]byte(v), &d)
	if err != nil {
		return err
	}
	if d == "true" {
		*t = true
	} else {
		*t = false
	}
	return nil
}

func (t StringBool) Value() (driver.Value, error) {

	var by []byte

	if t {
		by, _ = json.Marshal(true)
	} else {
		by, _ = json.Marshal(false)
	}

	return string(by), nil
}

type Skills [][]*string

func (t *Skills) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t Skills) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}

type TransName map[string]string

func (t *TransName) Scan(value interface{}) error {
	v, _ := value.(string)

	return json.Unmarshal([]byte(v), t)
}

func (t TransName) Value() (driver.Value, error) {
	by, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(by), nil
}
