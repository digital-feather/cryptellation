package exchanges

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/period"
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
