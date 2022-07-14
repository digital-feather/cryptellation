package event

import (
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/candlesticks"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/candlestick"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/tick"
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
	c *candlesticks.Candlestick,
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
