package service

import (
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

const MockExchangeName = "mock_exchange"

type MockExchangeService struct {
}

func (mes *MockExchangeService) ListenSymbol(symbol string) (chan tick.Tick, chan struct{}, error) {
	tickChan := make(chan tick.Tick, 200)

	for i := int64(0); i < 100; i++ {
		tickChan <- tick.Tick{
			Time:       time.Unix(i, 0),
			PairSymbol: symbol,
			Exchange:   MockExchangeName,
			Price:      float64(100 + i),
		}
	}

	return tickChan, make(chan struct{}), nil
}
