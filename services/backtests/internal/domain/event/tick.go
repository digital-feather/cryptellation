package event

import (
	"encoding/json"
	"time"

	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/tick"
)

type TickEvent struct {
	Event
	Content tick.Tick
}

func NewTickEvent(ts time.Time, t tick.Tick) TickEvent {
	return TickEvent{
		Event: Event{
			Type: TypeIsTick,
			Time: ts,
		},
		Content: t,
	}
}

func (t TickEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TickEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
