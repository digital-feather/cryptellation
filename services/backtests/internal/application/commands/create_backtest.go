package commands

import (
	"context"

	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/backtest"
	"golang.org/x/xerrors"
)

type CreateBacktestHandler struct {
	repository vdb.Port
	pubsub     pubsub.Port
}

func NewCreateBacktestHandler(repository vdb.Port, ps pubsub.Port) CreateBacktestHandler {
	if repository == nil {
		panic("nil repository")
	}

	if ps == nil {
		panic("nil pubsub")
	}

	return CreateBacktestHandler{
		repository: repository,
		pubsub:     ps,
	}
}

func (h CreateBacktestHandler) Handle(ctx context.Context, req backtest.NewPayload) (id uint, err error) {
	bt, err := backtest.New(ctx, req)
	if err != nil {
		return 0, xerrors.Errorf("creating a new backtest from request: %w", err)
	}

	err = h.repository.CreateBacktest(ctx, &bt)
	if err != nil {
		return 0, xerrors.Errorf("adding backtest to vdb: %w", err)
	}

	return bt.ID, nil
}
