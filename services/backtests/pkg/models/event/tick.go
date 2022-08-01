package event

import (
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	candlesticksProto "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

func NewTickEvent(t time.Time, content tick.Tick) Event {
	return Event{
		Type:    TypeIsTick,
		Time:    t,
		Content: content,
	}
}

func TickEventFromCandlestick(
	exchange, pairSymbol string,
	currentPriceType candlestick.PriceType,
	c *candlesticksProto.Candlestick,
) (Event, error) {
	t, err := time.Parse(time.RFC3339, c.Time)
	if err != nil {
		return Event{}, fmt.Errorf("error when parsing time from candlestick: %w", err)
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
