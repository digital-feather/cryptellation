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

func (t Type) String() string {
	return string(t)
}

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

type Interface interface {
	GetType() Type
	GetTime() time.Time
	GetContent() interface{}
}

func OnlyKeepEarliestSameTimeEvents(evts []Interface, endTime time.Time) (earliestTime time.Time, filtered []Interface) {
	earliestTime = endTime
	for _, e := range evts {
		t := e.GetTime()

		if earliestTime.After(t) {
			earliestTime = t
			filtered = make([]Interface, 0, len(evts))
		}

		if earliestTime.Equal(t) {
			filtered = append(filtered, e)
		}
	}

	return earliestTime, filtered
}
