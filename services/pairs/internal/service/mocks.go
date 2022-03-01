package service

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
)

type AssetsServiceMock struct {
}

func (asm AssetsServiceMock) ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error) {
	return []asset.Asset{
		{
			Symbol: "ETH",
		},
		{
			Symbol: "USDC",
		},
	}, nil
}
