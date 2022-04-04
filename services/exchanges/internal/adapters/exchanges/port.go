package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
)

type Port interface {
	Infos(ctx context.Context) (exchange.Exchange, error)
}
