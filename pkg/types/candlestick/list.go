package candlestick

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/period"
	"github.com/cryptellation/cryptellation/pkg/types/timeserie"
)

var (
	ErrPeriodMismatch   = errors.New("period-mismatch")
	ErrCandlestickType  = errors.New("struct-not-candlestick")
	ErrExchangeMismatch = errors.New("exchange-mismatch")
	ErrPairMismatch     = errors.New("pair-mismatch")
)

type ListID struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
}

type List struct {
	id           ListID
	candleSticks *timeserie.TimeSerie
}

func NewList(id ListID) *List {
	return &List{
		id:           id,
		candleSticks: timeserie.New(),
	}
}

func (l List) ID() ListID {
	return l.id
}

func (l List) ExchangeName() string {
	return l.id.ExchangeName
}

func (l List) PairSymbol() string {
	return l.id.PairSymbol
}

func (l List) Period() period.Symbol {
	return l.id.Period
}

func (l List) Len() int {
	return l.candleSticks.Len()
}

func (l List) Get(t time.Time) (Candlestick, bool) {
	data, exist := l.candleSticks.Get(t)
	if !exist {
		return Candlestick{}, false
	}
	return data.(Candlestick), true
}

func (l *List) Set(t time.Time, c Candlestick) error {
	if !l.id.Period.IsAligned(t) {
		return ErrPeriodMismatch
	}

	l.candleSticks.Set(t, c)
	return nil
}

func (l *List) MergeTimeSeries(ts timeserie.TimeSerie, options *timeserie.MergeOptions) error {
	if err := ts.Loop(func(t time.Time, obj interface{}) (bool, error) {
		if _, isCandlestick := obj.(Candlestick); !isCandlestick {
			return false, ErrCandlestickType
		}
		return false, nil
	}); err != nil {
		return err
	}

	return l.candleSticks.Merge(ts, options)
}

func (l *List) Merge(l2 List, options *timeserie.MergeOptions) error {
	if l.id.ExchangeName != l2.id.ExchangeName {
		return ErrExchangeMismatch
	} else if l.id.PairSymbol != l2.id.PairSymbol {
		return ErrPairMismatch
	} else if l.id.Period != l2.id.Period {
		return ErrPeriodMismatch
	}

	return l.candleSticks.Merge(*l2.candleSticks, options)
}

func (l *List) ReplaceUncomplete(l2 List) error {
	return l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			ucs, exists := l2.Get(t)
			if exists {
				return false, l.Set(t, ucs)
			}
		}
		return false, nil
	})
}

func (l *List) HasUncomplete() bool {
	hasUncomplete := false

	l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if cs.Uncomplete {
			hasUncomplete = true
			return true, nil
		}
		return false, nil
	})

	return hasUncomplete
}

func (l *List) Delete(t ...time.Time) {
	l.candleSticks.Delete(t...)
}

func (l *List) Loop(callback func(t time.Time, cs Candlestick) (bool, error)) error {
	return l.candleSticks.Loop(func(t time.Time, obj interface{}) (bool, error) {
		cs := obj.(Candlestick)
		return callback(t, cs)
	})
}

func (l List) First() (time.Time, Candlestick, bool) {
	t, data, ok := l.candleSticks.First()
	if !ok {
		return t, Candlestick{}, false
	}
	return t, data.(Candlestick), true
}

func (l List) Last() (time.Time, Candlestick, bool) {
	t, data, ok := l.candleSticks.Last()
	if !ok {
		return t, Candlestick{}, false
	}
	return t, data.(Candlestick), true
}

func (l List) Extract(start, end time.Time) *List {
	el := NewList(l.id)
	el.candleSticks = l.candleSticks.Extract(start, end)
	return el
}

func (l List) FirstN(limit uint) *List {
	el := NewList(l.id)
	el.candleSticks = l.candleSticks.FirstN(limit)
	return el
}

type ExportedList struct {
	Exchange     string                    `json:"exchange"`
	Pair         string                    `json:"pair"`
	Period       period.Symbol             `json:"period"`
	Candlesticks map[time.Time]Candlestick `json:"candlesticks,omitempty"`
}

func (l List) MarshalJSON() ([]byte, error) {
	var el ExportedList

	el.Exchange = l.ExchangeName()
	el.Pair = l.PairSymbol()
	el.Period = l.Period()

	el.Candlesticks = make(map[time.Time]Candlestick)
	if err := l.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		el.Candlesticks[t] = cs
		return false, nil
	}); err != nil {
		return nil, err
	}

	return json.Marshal(&el)
}

func (l *List) UnmarshalJSON(data []byte) error {
	var el ExportedList
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	per := el.Period
	if err := el.Period.Validate(); err != nil {
		return err
	}

	l.id.ExchangeName = el.Exchange
	l.id.PairSymbol = el.Pair
	l.id.Period = per
	l.candleSticks = timeserie.New()
	for t, c := range el.Candlesticks {
		if err := l.Set(t, c); err != nil {
			return fmt.Errorf("[candlestick %s]:%s", t.Format(time.RFC3339Nano), err)
		}
	}

	return nil
}

func MergeIntoOne(csl *List, per period.Symbol) (time.Time, Candlestick) {
	if csl.Len() == 0 {
		return time.Unix(0, 0), Candlestick{}
	}

	tsFirst, mcs, _ := csl.First()
	mts := per.RoundTime(tsFirst)

	csl.Loop(func(t time.Time, cs Candlestick) (bool, error) {
		if !per.RoundTime(t).Equal(mts) {
			return true, nil
		}

		if t.Equal(tsFirst) {
			return false, nil
		}

		if cs.High > mcs.High {
			mcs.High = cs.High
		}
		if cs.Low < mcs.Low {
			mcs.Low = cs.Low
		}
		mcs.Volume += cs.Volume
		mcs.Close = cs.Close

		return false, nil
	})

	return mts, mcs
}
