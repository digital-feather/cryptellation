package assets

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
)

type Port interface {
	ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error)
}
