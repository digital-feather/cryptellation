package vdb

import (
	"context"
	"time"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type Port interface {
	IncrementSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	DecrementSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	GetSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	ClearSymbolListenersCount(ctx context.Context) error
}
