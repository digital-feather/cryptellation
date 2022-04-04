package queriesBacktest

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/order"
)

type GetOrders struct {
	repository vdb.Port
}

func NewGetOrders(repository vdb.Port) GetOrders {
	if repository == nil {
		panic("nil repository")
	}

	return GetOrders{
		repository: repository,
	}
}

func (h GetOrders) Handle(ctx context.Context, backtestId uint) ([]order.Order, error) {
	bt, err := h.repository.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}
