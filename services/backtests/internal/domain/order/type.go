package order

import "errors"

var (
	ErrInvalidType = errors.New("invalid-order-type")
)

type Type string

const (
	TypeIsMarket Type = "market"
)

var (
	Types = []Type{
		TypeIsMarket,
	}
)

func (t Type) Validate() error {
	for _, vt := range Types {
		if t == vt {
			return nil
		}
	}

	return ErrInvalidType
}

func (t Type) String() string {
	return string(t)
}
