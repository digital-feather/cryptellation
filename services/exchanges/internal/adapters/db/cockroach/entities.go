package cockroach

import (
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/exchange"
)

type Pair struct {
	Symbol string `gorm:"primaryKey;autoIncrement:false"`
}

type Period struct {
	Symbol string `gorm:"primaryKey;autoIncrement:false"`
}

type Exchange struct {
	Name         string   `gorm:"primaryKey;autoIncrement:false"`
	Pairs        []Pair   `gorm:"many2many:exchanges_pairs;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Periods      []Period `gorm:"many2many:exchanges_periods;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Fees         float64
	LastSyncTime time.Time
}

func (e *Exchange) FromModel(model exchange.Exchange) {
	e.Name = model.Name

	e.Pairs = make([]Pair, len(model.Pairs))
	for i, p := range model.Pairs {
		e.Pairs[i] = Pair{
			Symbol: p,
		}
	}

	e.Periods = make([]Period, len(model.Periods))
	for i, p := range model.Periods {
		e.Periods[i] = Period{
			Symbol: p,
		}
	}

	e.Fees = model.Fees
	e.LastSyncTime = model.LastSyncTime
}

func (e Exchange) ToModel() exchange.Exchange {
	m := exchange.Exchange{
		Name:         e.Name,
		Pairs:        make([]string, len(e.Pairs)),
		Periods:      make([]string, len(e.Periods)),
		Fees:         e.Fees,
		LastSyncTime: e.LastSyncTime,
	}

	for i, p := range e.Pairs {
		m.Pairs[i] = p.Symbol
	}

	for i, p := range e.Periods {
		m.Periods[i] = p.Symbol
	}

	return m
}
