package db

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/candlestick"
)

type Port interface {
	CreateCandlesticks(ctx context.Context, cs *candlestick.List) error
	ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error
	UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error
	DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error
}
