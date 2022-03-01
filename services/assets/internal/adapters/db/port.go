package db

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
)

type Port interface {
	CreateAssets(ctx context.Context, assets ...asset.Asset) error
	ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error)
	UpdateAssets(ctx context.Context, assets ...asset.Asset) error
	DeleteAssets(ctx context.Context, assets ...asset.Asset) error
}