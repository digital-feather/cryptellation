package cockroach

import "github.com/cryptellation/cryptellation/services/assets/internal/domain/asset"

type Asset struct {
	Symbol string `gorm:"primaryKey"`
}

func (a *Asset) FromModel(model asset.Asset) {
	a.Symbol = model.Symbol
}

func (a *Asset) ToModel() asset.Asset {
	return asset.Asset{
		Symbol: a.Symbol,
	}
}
