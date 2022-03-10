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

type Candlestick struct {
	Open       float64 `bson:"open"     json:"open,omitempty"`
	High       float64 `bson:"high"     json:"high,omitempty"`
	Low        float64 `bson:"low"      json:"low,omitempty"`
	Close      float64 `bson:"close"    json:"close,omitempty"`
	Volume     float64 `bson:"volume"   json:"volume,omitempty"`
	Uncomplete bool    `bson:"complete" json:"uncomplete,omitempty"`
}

func (cs Candlestick) Equal(b Candlestick) bool {
	o := cs.Open == b.Open
	h := cs.High == b.High
	l := cs.Low == b.Low
	c := cs.Close == b.Close
	v := cs.Volume == b.Volume
	u := cs.Uncomplete == b.Uncomplete
	return o && h && l && c && v && u
}

func (cs Candlestick) PriceByType(pt PriceType) float64 {
	switch pt {
	case PriceTypeIsOpen:
		return cs.Open
	case PriceTypeIsHigh:
		return cs.High
	case PriceTypeIsLow:
		return cs.Low
	default:
		return cs.Close
	}
}
