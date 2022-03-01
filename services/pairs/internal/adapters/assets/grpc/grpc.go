package assetsgrpc

import (
	"context"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"golang.org/x/xerrors"
)

type AssetsGrpc struct {
	client assets.AssetsServiceClient
}

func NewAssetsGrpc(client assets.AssetsServiceClient) AssetsGrpc {
	return AssetsGrpc{client: client}
}

func (ag AssetsGrpc) ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error) {
	resp, err := ag.client.ReadAssets(ctx, &assets.ReadAssetsRequest{
		Symbols: symbols,
	})

	if err != nil {
		return nil, xerrors.Errorf("reading assets: %w", err)
	}

	assets := make([]asset.Asset, len(resp.Assets))
	for i, a := range resp.Assets {
		assets[i] = asset.Asset{
			Symbol: a.Symbol,
		}
	}

	return assets, nil
}
