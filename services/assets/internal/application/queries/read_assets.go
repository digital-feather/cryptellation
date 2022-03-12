package queries

import (
	"context"

	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/assets/pkg/asset"
	"golang.org/x/xerrors"
)

type ReadAssetsHandler struct {
	repository db.Port
}

func NewReadAssetsHandler(repository db.Port) ReadAssetsHandler {
	if repository == nil {
		panic("nil repository")
	}

	return ReadAssetsHandler{
		repository: repository,
	}
}

func (h ReadAssetsHandler) Handle(ctx context.Context, symbols []string) ([]asset.Asset, error) {
	as, err := h.repository.ReadAssets(ctx, symbols...)
	if err != nil {
		return as, xerrors.Errorf("handling assets reading: %w", err)
	}

	return as, nil
}
