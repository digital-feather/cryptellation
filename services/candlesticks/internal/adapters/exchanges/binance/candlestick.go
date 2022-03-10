package binance

import (
	"context"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/exchanges"
)

// CandlestickService is the real service for candlesticks
type CandlestickService struct {
	service    *client.KlinesService
	pairSymbol string
	period     period.Symbol
}

// Do will execute a request for candlesticks
func (s *CandlestickService) Do(ctx context.Context) (*candlestick.List, error) {
	// Get KLines
	kl, err := s.service.Do(ctx)
	if err != nil {
		return nil, err
	}

	// Change them to right format
	return KLinesToCandlesticks(s.pairSymbol, s.period, kl, time.Now())
}

// EndTime will specify the time where the list ends for next candlesticks request
func (s *CandlestickService) EndTime(endTime time.Time) exchanges.CandlesticksService {
	binanceTime := TimeToKLineTime(endTime)
	s.service.EndTime(binanceTime)
	return s
}

// StartTime will specify the time where the list starts for next candlesticks request
func (s *CandlestickService) StartTime(startTime time.Time) exchanges.CandlesticksService {
	binanceTime := TimeToKLineTime(startTime)
	s.service.StartTime(binanceTime)
	return s
}

// Limit will specify the number of candlesticks the list should have at its maximum
// If the limit is higher than the default limit, it will be limited to this one
func (s *CandlestickService) Limit(limit int) exchanges.CandlesticksService {
	s.service.Limit(limit)
	return s
}

func (s *CandlestickService) Period() period.Symbol {
	return s.period
}

func (s *CandlestickService) PairSymbol() string {
	return s.pairSymbol
}
