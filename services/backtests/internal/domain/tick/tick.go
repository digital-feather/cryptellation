package tick

import (
	"encoding/json"
)

type Tick struct {
	PairSymbol string
	Price      float64
	Exchange   string
}

func (t Tick) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tick) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
