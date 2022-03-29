package queriesBacktest

import (
	"context"

	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/account"
)

type GetAccounts struct {
	repository vdb.Port
}

func NewGetAccounts(repository vdb.Port) GetAccounts {
	if repository == nil {
		panic("nil repository")
	}

	return GetAccounts{
		repository: repository,
	}
}

func (h GetAccounts) Handle(ctx context.Context, backtestId uint) (map[string]account.Account, error) {
	bt, err := h.repository.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}
