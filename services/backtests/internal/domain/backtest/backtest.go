package backtest

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/account"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/candlestick"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/event"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/order"
	"github.com/cryptellation/cryptellation/services/candlesticks/pkg/period"
	"golang.org/x/xerrors"
)

var (
	ErrTickSubscriptionAlreadyExists = errors.New("tick-subscription-already-exists")
	ErrInvalidExchange               = errors.New("invalid-exchange")
	ErrNotEnoughAsset                = errors.New("no-enough-asset")
	ErrNoDataForOrderValidation      = errors.New("no-data-for-order-validation")
	ErrStartAfterEnd                 = errors.New("start-after-end")
)

// Current tick based on candlestick step
type CurrentCsTick struct {
	Time      time.Time
	PriceType candlestick.PriceType
}

type Backtest struct {
	ID                  uint
	StartTime           time.Time
	CurrentCsTick       CurrentCsTick
	EndTime             time.Time
	Accounts            map[string]account.Account
	PeriodBetweenEvents period.Symbol
	TickSubscribers     []event.Subscription
	Orders              []order.Order
}

type NewPayload struct {
	Accounts              map[string]account.Account
	StartTime             time.Time
	EndTime               *time.Time
	DurationBetweenEvents *time.Duration
}

func (payload *NewPayload) EmptyFieldsToDefault() *NewPayload {
	if payload.EndTime == nil {
		payload.EndTime = defaultEndTime()
	}

	if payload.DurationBetweenEvents == nil {
		d := time.Minute
		payload.DurationBetweenEvents = &d
	}

	return payload
}

func (payload NewPayload) Validate() error {
	if !payload.StartTime.Before(*payload.EndTime) {
		return ErrStartAfterEnd
	}

	for exchangeName, a := range payload.Accounts {
		if exchangeName == "" {
			return ErrInvalidExchange
		}

		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func defaultEndTime() *time.Time {
	t := time.Now()
	return &t
}

func New(ctx context.Context, payload NewPayload) (Backtest, error) {
	if err := payload.EmptyFieldsToDefault().Validate(); err != nil {
		return Backtest{}, err
	}

	per, err := period.FromDuration(*payload.DurationBetweenEvents)
	if err != nil {
		return Backtest{}, xerrors.Errorf("invalid duration between events: %w", err)
	}

	return Backtest{
		StartTime: payload.StartTime,
		CurrentCsTick: CurrentCsTick{
			Time:      payload.StartTime,
			PriceType: candlestick.PriceTypeIsOpen,
		},
		EndTime:             *payload.EndTime,
		Accounts:            payload.Accounts,
		PeriodBetweenEvents: per,
		TickSubscribers:     make([]event.Subscription, 0),
		Orders:              make([]order.Order, 0),
	}, nil
}

func (bt Backtest) MarshalBinary() ([]byte, error) {
	return json.Marshal(bt)
}

func (bt *Backtest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, bt)
}

func (bt *Backtest) Advance() {
	bt.advanceThroughTicks()
}

func (bt *Backtest) advanceThroughTicks() {
	switch bt.CurrentCsTick.PriceType {
	case candlestick.PriceTypeIsOpen:
		bt.CurrentCsTick.PriceType = candlestick.PriceTypeIsHigh
	case candlestick.PriceTypeIsHigh:
		bt.CurrentCsTick.PriceType = candlestick.PriceTypeIsLow
	case candlestick.PriceTypeIsLow:
		bt.CurrentCsTick.PriceType = candlestick.PriceTypeIsClose
	case candlestick.PriceTypeIsClose:
		bt.SetCurrentTime(bt.CurrentCsTick.Time.Add(bt.PeriodBetweenEvents.Duration()))
	default:
		bt.CurrentCsTick.PriceType = candlestick.PriceTypeIsOpen
	}
}

func (bt *Backtest) SetCurrentTime(ts time.Time) {
	bt.CurrentCsTick.Time = ts
	bt.CurrentCsTick.PriceType = candlestick.PriceTypeIsOpen
}

func (bt Backtest) Done() bool {
	return !bt.CurrentCsTick.Time.Before(bt.EndTime)
}

func (bt *Backtest) CreateTickSubscription(exchangeName string, pairSymbol string) (event.Subscription, error) {
	for _, ts := range bt.TickSubscribers {
		if ts.ExchangeName == exchangeName && ts.PairSymbol == pairSymbol {
			return event.Subscription{}, ErrTickSubscriptionAlreadyExists
		}
	}

	s := event.Subscription{
		ID:           len(bt.TickSubscribers),
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
	}
	bt.TickSubscribers = append(bt.TickSubscribers, s)

	return s, nil
}

// func (bt *Backtest) AddOrder(ctx context.Context, ord *order.Order) error {
// 	exchangeAccount, ok := bt.Accounts[ord.ExchangeName]
// 	if !ok {
// 		return ErrInvalidExchange
// 	}

// 	csl := candlestick.NewList(candlestick.ListID{
// 		Exchange: ord.ExchangeName,
// 		Pair:     ord.Pair,
// 		Period:   bt.TimeBetweenEvent,
// 	})

// 	if err := historian.FromContext(ctx).Candlesticks().Get(csl, bt.CurrentCsTick.Time, bt.CurrentCsTick.Time, 0); err != nil {
// 		return err
// 	}

// 	cs, ok := csl.Get(bt.CurrentCsTick.Time)
// 	if !ok {
// 		return ErrNoDataForOrderValidation
// 	}

// 	var price float64
// 	switch bt.CurrentCsTick.PriceType {
// 	case candlestick.PriceTypeIsOpen:
// 		price = cs.Open
// 	case candlestick.PriceTypeIsHigh:
// 		price = cs.High
// 	case candlestick.PriceTypeIsLow:
// 		price = cs.Low
// 	case candlestick.PriceTypeIsClose:
// 		fallthrough
// 	default:
// 		price = cs.Close
// 	}

// 	quoteEquivalentQty := price * ord.Quantity
// 	if ord.Side == order.SideIsBuy {
// 		available, ok := exchangeAccount.Balances[ord.Pair.QuoteSymbol]
// 		if !ok {
// 			return ErrNotEnoughAsset
// 		} else if quoteEquivalentQty > available {
// 			return ErrNotEnoughAsset
// 		}

// 		bt.Accounts[ord.ExchangeName].Balances[ord.Pair.QuoteSymbol] -= quoteEquivalentQty
// 		bt.Accounts[ord.ExchangeName].Balances[ord.Pair.BaseSymbol] += ord.Quantity
// 	} else {
// 		available, ok := exchangeAccount.Balances[ord.Pair.BaseSymbol]
// 		if !ok {
// 			return ErrNotEnoughAsset
// 		} else if ord.Quantity > available {
// 			return ErrNotEnoughAsset
// 		}

// 		bt.Accounts[ord.ExchangeName].Balances[ord.Pair.QuoteSymbol] += quoteEquivalentQty
// 		bt.Accounts[ord.ExchangeName].Balances[ord.Pair.BaseSymbol] -= ord.Quantity
// 	}

// 	ord.ID = uint(len(bt.Orders))
// 	ord.Status = order.StatusIsFilled

// 	bt.Orders = append(bt.Orders, *ord)
// 	return nil
// }
