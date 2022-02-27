package queries

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db"
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
	return h.repository.ReadAssets(ctx, symbols...)
}
