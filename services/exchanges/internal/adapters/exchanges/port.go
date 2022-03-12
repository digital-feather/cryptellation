package exchanges

import (
	"context"

	"github.com/cryptellation/cryptellation/services/exchanges/pkg/exchange"
)

type Port interface {
	Infos(ctx context.Context) (exchange.Exchange, error)
}
