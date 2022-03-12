package service

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/services/exchanges/internal/domain/exchange"
)

type MockExchangeService struct {
}

func (mes MockExchangeService) Infos(ctx context.Context) (exchange.Exchange, error) {
	return exchange.Exchange{
		Name:           "mock_exchange",
		PairsSymbols:   []string{"ABC-DEF", "IJK-LMN"},
		PeriodsSymbols: []string{"M1", "M3"},
		Fees:           0.1,
		LastSyncTime:   time.Now().UTC(),
	}, nil
}
