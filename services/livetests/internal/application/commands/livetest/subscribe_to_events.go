package cmdLivetest

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
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

func (h SubscribeToEventsHandler) Handle(ctx context.Context, livetestId uint, exchange, pairSymbol string) error {
	return h.repository.LockedLivetest(livetestId, func() error {
		bt, err := h.repository.ReadLivetest(ctx, livetestId)
		if err != nil {
			return fmt.Errorf("cannot get livetest: %w", err)
		}

		if _, err = bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		if err := h.repository.UpdateLivetest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update livetest: %w", err)
		}

		return nil
	})
}
