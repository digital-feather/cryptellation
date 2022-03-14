package event

import (
	"encoding/json"
	"time"
)

type EndEvent struct {
	Event
}

func NewEndEvent(t time.Time) EndEvent {
	return EndEvent{
		Event: Event{
			Type: TypeIsEnd,
			Time: t,
		},
	}
}

func (e EndEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EndEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}
