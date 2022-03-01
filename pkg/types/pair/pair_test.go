package pair

import "testing"

func TestNew(t *testing.T) {
	pn := New("ETH", "USDT")
	p := Pair{BaseSymbol: "ETH", QuoteSymbol: "USDT"}

	if p != pn {
		t.Error("Should be equal")
	}
}

func TestPairString(t *testing.T) {
	p := Pair{BaseSymbol: "ETH", QuoteSymbol: "USDT"}

	if p.String() != "ETH-USDT" {
		t.Error("Wrong string format")
	}
}

func TestUniqueArray(t *testing.T) {
	p1 := []Pair{New("ABC", "DEF"), New("DEF", "XYZ")}
	p2 := []Pair{New("ABC", "DEF"), New("ABC", "XYZ")}
	p3 := []Pair{New("ABC", "DEF"), New("ABC", "XYZ"), New("DEF", "XYZ")}

	m := UniqueArray(p2, p1)
	if len(m) != 3 || m[0] != p3[0] || m[1] != p3[1] || m[2] != p3[2] {
		t.Error(p3, m)
	}
}
