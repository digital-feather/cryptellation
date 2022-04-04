package cmdBacktest

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"golang.org/x/xerrors"
)

type CreateHandler struct {
	repository vdb.Port
}

func NewCreateHandler(repository vdb.Port) CreateHandler {
	if repository == nil {
		panic("nil repository")
	}

	return CreateHandler{
		repository: repository,
	}
}

func (h CreateHandler) Handle(ctx context.Context, req backtest.NewPayload) (id uint, err error) {
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
