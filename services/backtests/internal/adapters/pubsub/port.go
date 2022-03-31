package pubsub

import (
	"context"

	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/event"
)

type Port interface {
	Publish(ctx context.Context, backtestID uint, event event.Interface) error
	Subscribe(ctx context.Context, backtestID uint) (Subscriber, error)
}

type Subscriber interface {
	Channel() <-chan event.Interface
	Close() error
}