package service

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/exchanges"
)

type MockExchangeService struct {
}

func (m *MockExchangeService) Candlesticks(pairSymbol string, period period.Symbol) (exchanges.CandlesticksService, error) {
	return &CandlesticksService{
		pairSymbol: pairSymbol,
		period:     period,
	}, nil
}

type CandlesticksService struct {
	pairSymbol string
	period     period.Symbol
}

func (m *CandlesticksService) Do(ctx context.Context) (*candlestick.List, error) {
	return candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   m.PairSymbol(),
		Period:       m.Period(),
	}), nil
}

func (m *CandlesticksService) StartTime(startTime time.Time) exchanges.CandlesticksService {
	return m
}

func (m *CandlesticksService) EndTime(endTime time.Time) exchanges.CandlesticksService {
	return m
}

func (m *CandlesticksService) Limit(limit int) exchanges.CandlesticksService {
	return m
}

func (m CandlesticksService) Period() period.Symbol {
	return m.period
}

func (m CandlesticksService) PairSymbol() string {
	return m.pairSymbol
}
