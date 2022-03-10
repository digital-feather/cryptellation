package exchanges

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
)

type Port interface {
	Candlesticks(pairSymbol string, period period.Symbol) (CandlesticksService, error)
}

type CandlesticksService interface {
	Do(ctx context.Context) (*candlestick.List, error)
	StartTime(startTime time.Time) CandlesticksService
	EndTime(endTime time.Time) CandlesticksService
	Limit(limit int) CandlesticksService
	Period() period.Symbol
	PairSymbol() string
}
