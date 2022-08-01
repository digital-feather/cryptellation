package event

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    Type
	Time    time.Time
	Content interface{}
}

func (e Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

func OnlyKeepEarliestSameTimeEvents(evts []Event, endTime time.Time) (earliestTime time.Time, filtered []Event) {
	earliestTime = endTime
	for _, e := range evts {
		if earliestTime.After(e.Time) {
			earliestTime = e.Time
			filtered = make([]Event, 0, len(evts))
		}

		if earliestTime.Equal(e.Time) {
			filtered = append(filtered, e)
		}
	}

	return earliestTime, filtered
}
