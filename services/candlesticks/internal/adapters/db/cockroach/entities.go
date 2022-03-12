package cockroach

import (
	"fmt"
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/candlestick"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/period"
	"golang.org/x/xerrors"
)

type Candlestick struct {
	ExchangeName string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PairSymbol   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PeriodSymbol string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Time         time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Open         float64
	High         float64
	Low          float64
	Close        float64
	Volume       float64
	Uncomplete   bool
}

func (c *Candlestick) FromModel(exchange, pair, period string, t time.Time, model candlestick.Candlestick) {
	c.ExchangeName = exchange
	c.PairSymbol = pair
	c.PeriodSymbol = period
	c.Time = t

	c.Open = model.Open
	c.High = model.High
	c.Low = model.Low
	c.Close = model.Close
	c.Volume = model.Volume
	c.Uncomplete = model.Uncomplete
}

func (c Candlestick) ToModel() (exchange, pair, period string, t time.Time, model candlestick.Candlestick) {
	return c.ExchangeName,
		c.PairSymbol,
		c.PeriodSymbol,
		c.Time,
		candlestick.Candlestick{
			Open:       c.Open,
			High:       c.High,
			Low:        c.Low,
			Close:      c.Close,
			Volume:     c.Volume,
			Uncomplete: c.Uncomplete,
		}
}

func FromModelListToEntityList(list *candlestick.List) []Candlestick {
	entities := make([]Candlestick, 0, list.Len())
	_ = list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		c := Candlestick{}
		c.FromModel(list.ExchangeName(), list.PairSymbol(), list.Period().String(), t, cs)
		entities = append(entities, c)
		return false, nil
	})

	return entities
}

func FromEntityListToModelList(entities []Candlestick) (*candlestick.List, error) {
	if entities == nil || len(entities) == 0 {
		return nil, xerrors.New("no entities provided")
	}

	periodSymbol := entities[0].PeriodSymbol
	per, err := period.FromString(periodSymbol)
	if err != nil {
		return nil, xerrors.Errorf("from candlestick entity to list: %w", err)
	}

	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: entities[0].ExchangeName,
		PairSymbol:   entities[0].PairSymbol,
		Period:       per,
	})

	for _, e := range entities {
		if e.ExchangeName != list.ExchangeName() {
			txt := fmt.Sprintf("incompatible exchanges for same list: %q and %q", e.ExchangeName, list.ExchangeName())
			return nil, xerrors.New(txt)
		}

		if e.PairSymbol != list.PairSymbol() {
			txt := fmt.Sprintf("incompatible pair for same list: %q and %q", e.PairSymbol, list.PairSymbol())
			return nil, xerrors.New(txt)
		}

		if e.PeriodSymbol != list.Period().String() {
			txt := fmt.Sprintf("incompatible period for same list: %q and %q", e.PeriodSymbol, list.Period().String())
			return nil, xerrors.New(txt)
		}

		list.Set(e.Time, candlestick.Candlestick{
			Open:       e.Open,
			High:       e.High,
			Low:        e.Low,
			Close:      e.Close,
			Volume:     e.Volume,
			Uncomplete: e.Uncomplete,
		})
	}

	return list, nil
}
