package queries

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db"
	"golang.org/x/xerrors"
)

type ReadPairsHandler struct {
	repository db.Port
}

func NewReadPairsHandler(repository db.Port) ReadPairsHandler {
	if repository == nil {
		panic("nil repository")
	}

	return ReadPairsHandler{
		repository: repository,
	}
}

func (h ReadPairsHandler) Handle(ctx context.Context, pairs []string) ([]pair.Pair, error) {
	ps, err := h.repository.ReadPairs(ctx, pairs...)
	if err != nil {
		return ps, xerrors.Errorf("handling pairs reading: %w", err)
	}

	return ps, err
}
