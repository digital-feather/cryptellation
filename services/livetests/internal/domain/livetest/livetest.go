package livetest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/event"
)

var (
	ErrTickSubscriptionAlreadyExists = errors.New("tick-subscription-already-exists")
	ErrInvalidExchange               = errors.New("invalid-exchange")
)

type Livetest struct {
	ID              uint
	Accounts        map[string]account.Account
	TickSubscribers []event.Subscription
}

type NewPayload struct {
	Accounts map[string]account.Account
}

func (payload *NewPayload) EmptyFieldsToDefault() *NewPayload {
	return payload
}

func (payload NewPayload) Validate() error {
	for exchangeName, a := range payload.Accounts {
		if exchangeName == "" {
			return fmt.Errorf("error with exchange %q in new livetest payload: %w", exchangeName, ErrInvalidExchange)
		}

		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func New(ctx context.Context, payload NewPayload) (Livetest, error) {
	if err := payload.EmptyFieldsToDefault().Validate(); err != nil {
		return Livetest{}, err
	}

	return Livetest{
		Accounts:        payload.Accounts,
		TickSubscribers: make([]event.Subscription, 0),
	}, nil
}

func (bt Livetest) MarshalBinary() ([]byte, error) {
	return json.Marshal(bt)
}

func (bt *Livetest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, bt)
}

func (bt *Livetest) CreateTickSubscription(exchangeName string, pairSymbol string) (event.Subscription, error) {
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
