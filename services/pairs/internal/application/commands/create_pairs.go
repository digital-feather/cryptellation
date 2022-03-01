package commands

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/assets"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/pairs/internal/domain"
	"golang.org/x/xerrors"
)

type CreatePairsHandler struct {
	repository    db.Port
	assetsService assets.Port
}

func NewCreatePairsHandler(repository db.Port, assetsService assets.Port) CreatePairsHandler {
	if repository == nil {
		panic("nil repository")
	}

	if assetsService == nil {
		panic("nil assets service")
	}

	return CreatePairsHandler{
		repository:    repository,
		assetsService: assetsService,
	}
}

func (h CreatePairsHandler) Handle(ctx context.Context, pairs []pair.Pair) error {
	assets, err := h.assetsService.ReadAssets(ctx, getUniqueSliceOfSymbols(pairs)...)
	if err != nil {
		return xerrors.Errorf("reading assets: %w", err)
	}

	newPairs, err := domain.CreatePairs(pairs, assets)
	if err != nil {
		return xerrors.Errorf("creating pairs: %w", err)
	}

	err = h.repository.CreatePairs(ctx, newPairs...)
	if err != nil {
		return xerrors.Errorf("adding pairs to db: %w", err)
	}

	return nil
}

func getUniqueSliceOfSymbols(pairs []pair.Pair) []string {
	symbols := make(map[string]bool)
	for _, p := range pairs {
		symbols[p.BaseSymbol] = true
		symbols[p.QuoteSymbol] = true
	}

	symbolsList := make([]string, 0, len(symbols))
	for s := range symbols {
		symbolsList = append(symbolsList, s)
	}

	return symbolsList
}
