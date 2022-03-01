package domain

import (
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"golang.org/x/xerrors"
)

func CreatePairs(pairs []pair.Pair, existantAssets []asset.Asset) ([]pair.Pair, error) {
	assetsMap := make(map[string]asset.Asset)
	for _, a := range existantAssets {
		assetsMap[a.Symbol] = a
	}

	for _, p := range pairs {
		if _, ok := assetsMap[p.BaseAssetSymbol]; !ok {
			return nil, xerrors.Errorf("base symbol doesn't exist for pair %s", p)
		}

		if _, ok := assetsMap[p.QuoteAssetSymbol]; !ok {
			return nil, xerrors.Errorf("quote symbol doesn't exist for pair %s", p)
		}
	}

	return pairs, nil
}
