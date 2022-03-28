package tick

import (
	"encoding/json"
)

type Tick struct {
	PairSymbol string  `json:"pair_symbol"`
	Price      float64 `json:"price"`
	Exchange   string  `json:"exchange"`
}

func (t Tick) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tick) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
