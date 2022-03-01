package pair

import (
	"fmt"
)

type Pair struct {
	BaseSymbol  string `json:"base_symbol"`
	QuoteSymbol string `json:"quote_symbol"`
}

func New(baseSymbol, quoteSymbol string) Pair {
	return Pair{
		BaseSymbol:  baseSymbol,
		QuoteSymbol: quoteSymbol,
	}
}

func (p Pair) String() string {
	return fmt.Sprintf("%s-%s", p.BaseSymbol, p.QuoteSymbol)
}

func UniqueArray(pair1, pair2 []Pair) []Pair {
	tmp := make([]Pair, len(pair1))
	copy(tmp, pair1)

	for _, p2 := range pair2 {
		present := false
		for _, p1 := range pair1 {
			if p1 == p2 {
				present = true
				break
			}
		}

		if !present {
			tmp = append(tmp, p2)
		}
	}

	return tmp
}
