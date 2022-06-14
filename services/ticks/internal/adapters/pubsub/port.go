package pubsub

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

type Port interface {
	Publish(ctx context.Context, tick tick.Tick) error
	Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error)
}
