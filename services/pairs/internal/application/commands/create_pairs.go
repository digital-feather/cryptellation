package commands

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db"
	"golang.org/x/xerrors"
)

type CreatePairsHandler struct {
	repository db.Port
}

func NewCreatePairsHandler(repository db.Port) CreatePairsHandler {
	if repository == nil {
		panic("nil repository")
	}

	return CreatePairsHandler{
		repository: repository,
	}
}

func (h CreatePairsHandler) Handle(ctx context.Context, pairs []pair.Pair) error {
	err := h.repository.CreatePairs(ctx, pairs...)
	if err != nil {
		return xerrors.Errorf("adding pairs to db: %w", err)
	}

	return nil
}
