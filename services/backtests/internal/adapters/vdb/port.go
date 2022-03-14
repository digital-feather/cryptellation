package vdb

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/backtest"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type LockedBacktestCallback func() error

type Port interface {
	CreateBacktest(ctx context.Context, bt *backtest.Backtest) error
	ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error)
	UpdateBacktest(ctx context.Context, bt backtest.Backtest) error
	DeleteBacktest(ctx context.Context, bt backtest.Backtest) error

	LockedBacktest(id uint, fn LockedBacktestCallback) error
}
