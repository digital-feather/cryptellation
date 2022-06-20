package event

import (
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/status"
)

func NewStatusEvent(t time.Time, content status.Status) Event {
	return Event{
		Type:    TypeIsStatus,
		Time:    t,
		Content: content,
	}
}
