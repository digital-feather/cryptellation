package exchanges

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/exchange"
)

type Port interface {
	Infos(ctx context.Context) (exchange.Exchange, error)
}
