package event

import (
	"encoding/json"
	"time"
)

type Type string

const (
	TypeIsTick Type = "tick"
	TypeIsEnd  Type = "end"
)

type Event struct {
	Type Type
	Time time.Time
}

func (e Event) GetType() Type {
	return e.Type
}

func (e Event) GetTime() time.Time {
	return e.Time
}

func (e Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}
