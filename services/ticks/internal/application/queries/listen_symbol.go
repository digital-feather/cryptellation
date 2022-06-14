package queries

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

type ListenSymbolsHandler struct {
	pubsub pubsub.Port
}

func NewListenSymbolsHandler(ps pubsub.Port) ListenSymbolsHandler {
	if ps == nil {
		panic("nil pubsub")
	}

	return ListenSymbolsHandler{
		pubsub: ps,
	}
}

func (h ListenSymbolsHandler) Handle(ctx context.Context, exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return h.pubsub.Subscribe(ctx, pairSymbol)
}
