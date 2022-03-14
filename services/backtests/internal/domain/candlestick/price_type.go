package candlestick

import "errors"

var (
	ErrInvalidPriceType = errors.New("invalid-price-type")
)

type PriceType string

const (
	PriceTypeIsOpen  PriceType = "open"
	PriceTypeIsHigh  PriceType = "high"
	PriceTypeIsLow   PriceType = "low"
	PriceTypeIsClose PriceType = "close"
)

var PriceTypes = []PriceType{
	PriceTypeIsOpen,
	PriceTypeIsHigh,
	PriceTypeIsLow,
	PriceTypeIsClose,
}

func (pt PriceType) String() string {
	return string(pt)
}

func (pt PriceType) Validate() error {
	for _, vpt := range PriceTypes {
		if vpt.String() == pt.String() {
			return nil
		}
	}

	return ErrInvalidPriceType
}

// TODO add unmarshaling JSON with validation on pricetype

func PriceByType(open, high, low, close float64, pt PriceType) float64 {
	switch pt {
	case PriceTypeIsOpen:
		return open
	case PriceTypeIsHigh:
		return high
	case PriceTypeIsLow:
		return low
	default:
		return close
	}
}
