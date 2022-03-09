package service

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/exchange"
)

type MockExchangeService struct {
}

func (mes MockExchangeService) Infos(ctx context.Context) (exchange.Exchange, error) {
	return exchange.Exchange{
		Name:         "mock_exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC(),
	}, nil
}
