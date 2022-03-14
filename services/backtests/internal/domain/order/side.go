package order

import (
	"errors"
)

var (
	ErrInvalidSide = errors.New("invalid-order-side")
)

type Side string

const (
	SideIsBuy  Side = "buy"
	SideIsSell Side = "sell"
)

var (
	Sides = []Side{
		SideIsBuy,
		SideIsSell,
	}
)

func (s Side) Validate() error {
	for _, vs := range Sides {
		if s == vs {
			return nil
		}
	}

	return ErrInvalidSide
}

func (s Side) String() string {
	return string(s)
}
