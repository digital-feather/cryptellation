package pubsub

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
)

type Port interface {
	Publish(ctx context.Context, backtestID uint, event event.Event) error
	Subscribe(ctx context.Context, backtestID uint) (<-chan event.Event, error)
}
