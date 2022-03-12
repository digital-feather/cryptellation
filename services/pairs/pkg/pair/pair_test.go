package pair

import "testing"

func TestPairSymbol(t *testing.T) {
	p := Pair{BaseAssetSymbol: "ETH", QuoteAssetSymbol: "USDT"}

	if p.Symbol() != "ETH-USDT" {
		t.Error("Wrong symbol")
	}
}

func TestUniqueArray(t *testing.T) {
	p1 := []Pair{
		{BaseAssetSymbol: "ABC", QuoteAssetSymbol: "DEF"},
		{BaseAssetSymbol: "DEF", QuoteAssetSymbol: "XYZ"},
	}
	p2 := []Pair{
		{BaseAssetSymbol: "ABC", QuoteAssetSymbol: "DEF"},
		{BaseAssetSymbol: "ABC", QuoteAssetSymbol: "XYZ"},
	}
	p3 := []Pair{
		{BaseAssetSymbol: "ABC", QuoteAssetSymbol: "DEF"},
		{BaseAssetSymbol: "ABC", QuoteAssetSymbol: "XYZ"},
		{BaseAssetSymbol: "DEF", QuoteAssetSymbol: "XYZ"},
	}

	m := UniqueArray(p2, p1)
	if len(m) != 3 || m[0] != p3[0] || m[1] != p3[1] || m[2] != p3[2] {
		t.Error(p3, m)
	}
}
