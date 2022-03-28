package event

import (
	"encoding/json"
	"time"

	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/candlestick"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/tick"
	"golang.org/x/xerrors"
)

type TickEvent struct {
	Event
	Content tick.Tick
}

func NewTickEvent(t time.Time, ti tick.Tick) TickEvent {
	return TickEvent{
		Event: Event{
			Type: TypeIsTick,
			Time: t,
		},
		Content: ti,
	}
}

func (t TickEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TickEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t TickEvent) GetContent() interface{} {
	return t.Content
}

func TickEventFromCandlestick(
	exchange, pairSymbol string,
	currentPriceType candlestick.PriceType,
	c candlesticks.Candlestick,
) (TickEvent, error) {
	t, err := time.Parse(time.RFC3339, c.Time)
	if err != nil {
		return TickEvent{}, xerrors.Errorf("error when parsing time from candlestick: %w", err)
	}

	price := candlestick.PriceByType(
		float64(c.Open),
		float64(c.High),
		float64(c.Low),
		float64(c.Close),
		currentPriceType)

	return NewTickEvent(t, tick.Tick{
		PairSymbol: pairSymbol,
		Price:      price,
		Exchange:   exchange,
	}), nil
}
