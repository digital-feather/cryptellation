package cmdBacktest

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
)

type SubscribeToEventsHandler struct {
	repository vdb.Port
}

func NewSubscribeToEventsHandler(repository vdb.Port) SubscribeToEventsHandler {
	if repository == nil {
		panic("nil repository")
	}

	return SubscribeToEventsHandler{
		repository: repository,
	}
}

func (h SubscribeToEventsHandler) Handle(ctx context.Context, backtestId uint, exchange, pairSymbol string) error {
	return h.repository.LockedBacktest(backtestId, func() error {
		bt, err := h.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		if _, err = bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
