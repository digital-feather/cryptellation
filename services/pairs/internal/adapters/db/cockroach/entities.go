package cockroach

import "github.com/cryptellation/cryptellation/pkg/types/pair"

type Pair struct {
	BaseSymbol  string `gorm:"primaryKey;autoIncrement:false"`
	QuoteSymbol string `gorm:"primaryKey;autoIncrement:false"`
}

func (p *Pair) FromModel(model pair.Pair) {
	p.BaseSymbol = model.BaseSymbol
	p.QuoteSymbol = model.QuoteSymbol
}

func (p *Pair) ToModel() pair.Pair {
	return pair.New(p.BaseSymbol, p.QuoteSymbol)
}
