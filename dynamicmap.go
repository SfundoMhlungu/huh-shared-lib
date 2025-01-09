package main

import "encoding/json"

type Value struct {
	value string
}

func (v *Value) Get() interface{} {
	return v.value
}

func (v *Value) MarshalJSON() ([]byte, error) {

	return []byte(`"` + v.value + `"`), nil
}

func (v *Value) UnmarshalJSON(data []byte) error {

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	v.value = value
	return nil
}

type mapWithValues = map[string]*Value

func NewMapWithValues() mapWithValues {
	return make(mapWithValues)
}
