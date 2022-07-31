package candlestick

import "testing"

func TestCandlestickPriceTypes(t *testing.T) {
	validPriceTypes := []string{
		"open",
		"high",
		"low",
		"close",
	}

	for _, vpt := range validPriceTypes {
		if err := PriceType(vpt).Validate(); err != nil {
			t.Error("Price type should be valid: ", vpt)
		}
	}

	if err := PriceType("unknwon").Validate(); err == nil {
		t.Error("Price type should not be valid: unknwon")
	}
}

func TestCandlestickPriceByType(t *testing.T) {
	v := PriceByType(0, 1, 2, 3, PriceTypeIsOpen)
	if v != 0 {
		t.Error("Wrong value:", v)
	}

	v = PriceByType(0, 1, 2, 3, PriceTypeIsHigh)
	if v != 1 {
		t.Error("Wrong value:", v)
	}

	v = PriceByType(0, 1, 2, 3, PriceTypeIsLow)
	if v != 2 {
		t.Error("Wrong value:", v)
	}

	v = PriceByType(0, 1, 2, 3, PriceTypeIsClose)
	if v != 3 {
		t.Error("Wrong value:", v)
	}
}
