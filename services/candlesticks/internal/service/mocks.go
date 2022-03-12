package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/candlestick"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/period"
)

type MockExchangeService struct {
}

func (m *MockExchangeService) Candlesticks(pairSymbol string, period period.Symbol) (exchanges.CandlesticksService, error) {
	return &CandlesticksService{
		pairSymbol: pairSymbol,
		period:     period,
		start:      time.Unix(0, 0),
		end:        time.Now(),
	}, nil
}

type CandlesticksService struct {
	pairSymbol string
	period     period.Symbol
	start      time.Time
	end        time.Time
}

func (m *CandlesticksService) Do(ctx context.Context) (*candlestick.List, error) {
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   m.PairSymbol(),
		Period:       m.Period(),
	})

	fmt.Println("#", m.start, m.end)
	for i := m.start.Unix(); i < 60*1000; i++ {
		if time.Unix(i, 0).After(m.end) {
			break
		}

		cl.Set(time.Unix(i, 0), candlestick.Candlestick{
			Open:  float64(i),
			Close: 1234,
		})
	}

	return cl, nil
}

func (m *CandlesticksService) StartTime(startTime time.Time) exchanges.CandlesticksService {
	m.start = startTime
	return m
}

func (m *CandlesticksService) EndTime(endTime time.Time) exchanges.CandlesticksService {
	m.end = endTime
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
