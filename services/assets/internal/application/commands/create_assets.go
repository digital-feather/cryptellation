package commands

import (
	"context"

	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/assets/pkg/asset"
	"golang.org/x/xerrors"
)

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

func (h CreateAssetsHandler) Handle(ctx context.Context, assets []asset.Asset) error {
	err := h.repository.CreateAssets(ctx, assets...)
	if err != nil {
		return xerrors.Errorf("handling assets creation: %w", err)
	}

	return nil
}
