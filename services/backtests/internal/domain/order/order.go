package order

import "errors"

var (
	ErrInvalidOrderQty = errors.New("invalid-order-quantity")
)

type Order struct {
	ID           uint
	Type         Type
	ExchangeName string
	PairSymbol   string
	Side         Side
	Quantity     float64
	Status       Status
}

func (o Order) Validate() error {
	if err := o.Type.Validate(); err != nil {
		return err
	}

	if err := o.Side.Validate(); err != nil {
		return err
	}

	if o.Quantity <= 0 {
		return ErrInvalidOrderQty
	}

	return nil
}
