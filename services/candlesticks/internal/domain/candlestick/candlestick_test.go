package candlestick

import (
	"testing"
)

func TestCandlestickEqual(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 1000, false}

	if a.Equal(b) == false {
		t.Error("Candlesticks should be equal")
	}
}

func TestCandlestickNotEqualOpen(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{1, 1, 2, 3, 1000, false}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}

func TestCandlestickNotEqualHigh(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 2, 2, 3, 1000, false}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}

func TestCandlestickNotEqualLow(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 3, 3, 1000, false}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}

func TestCandlestickNotEqualClose(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 4, 1000, false}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}

func TestCandlestickNotEqualVolume(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 2000, false}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}

func TestCandlestickNotEqualUncomplete(t *testing.T) {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 1000, true}

	if a.Equal(b) {
		t.Error("Candlesticks should not be equal")
	}
}
