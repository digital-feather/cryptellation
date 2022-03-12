package db

import (
	"context"

	"github.com/cryptellation/cryptellation/services/exchanges/pkg/exchange"
)

type Port interface {
	CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error
	ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error)
	UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error
	DeleteExchanges(ctx context.Context, names ...string) error
}
