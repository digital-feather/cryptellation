package queriesBacktest

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
)

type ListenEventsHandler struct {
	pubsub pubsub.Port
}

func NewListenEventsHandler(ps pubsub.Port) ListenEventsHandler {
	if ps == nil {
		panic("nil pubsub")
	}

	return ListenEventsHandler{
		pubsub: ps,
	}
}

func (h ListenEventsHandler) Handle(ctx context.Context, backtestId uint64) (<-chan event.Event, error) {
	return h.pubsub.Subscribe(ctx, uint(backtestId))
}
