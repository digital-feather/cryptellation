package commands

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db"
)

type CreateAssets struct {
	Assets []asset.Asset
}

type CreateAssetsHandler struct {
	repository db.Port
}

func NewCreateAssetHandler(repository db.Port) CreateAssetsHandler {
	if repository == nil {
		panic("nil repository")
	}

	return CreateAssetsHandler{
		repository: repository,
	}
}

func (h CreateAssetsHandler) Handle(ctx context.Context, cmd CreateAssets) error {
	return h.repository.CreateAssets(ctx, cmd.Assets...)
}
