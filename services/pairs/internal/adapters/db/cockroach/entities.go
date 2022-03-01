package cockroach

import "github.com/cryptellation/cryptellation/pkg/types/pair"

type Pair struct {
	Symbol           string `gorm:"primaryKey;autoIncrement:false"`
	BaseAssetSymbol  string
	QuoteAssetSymbol string
}

func (p *Pair) FromModel(model pair.Pair) {
	p.Symbol = model.Symbol()
	p.BaseAssetSymbol = model.BaseAssetSymbol
	p.QuoteAssetSymbol = model.QuoteAssetSymbol
}

func (p *Pair) ToModel() pair.Pair {
	return pair.Pair{
		BaseAssetSymbol:  p.BaseAssetSymbol,
		QuoteAssetSymbol: p.QuoteAssetSymbol,
	}
}
