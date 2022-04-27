package commands

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
)

type UnregisterSymbolListenerHandler struct {
	vdb vdb.Port
}

func NewUnregisterSymbolListener(db vdb.Port) UnregisterSymbolListenerHandler {
	if db == nil {
		panic("nil vdb")
	}

	return UnregisterSymbolListenerHandler{
		vdb: db,
	}
}

func (h UnregisterSymbolListenerHandler) Handle(ctx context.Context, exchange, pairSymbol string) error {
	_, err := h.vdb.DecrementSymbolListenerCount(ctx, exchange, pairSymbol)
	if err != nil {
		return err
	}

	return nil
}
