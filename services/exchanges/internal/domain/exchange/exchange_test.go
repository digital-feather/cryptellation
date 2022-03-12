package exchange

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		Exchange1 Exchange
		Exchange2 Exchange
		Expected  Exchange
	}{
		{
			Exchange1: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Exchange2: Exchange{
				Name:           "exchange2",
				PairsSymbols:   []string{"ABC-DEF", "ABC-XYZ"},
				PeriodsSymbols: []string{"M1", "M3"},
				Fees:           0.2,
			},
			Expected: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ", "ABC-XYZ"},
				PeriodsSymbols: []string{"M1", "M15", "M3"},
				Fees:           0.1,
			},
		},
		{
			Exchange1: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Exchange2: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
			Expected: Exchange{
				Name:           "exchange1",
				PairsSymbols:   []string{"ABC-DEF", "DEF-XYZ"},
				PeriodsSymbols: []string{"M1", "M15"},
				Fees:           0.1,
			},
		},
	}

	for i, c := range cases {
		merged := c.Exchange1.Merge(c.Exchange2)
		if !reflect.DeepEqual(c.Expected, merged) {
			t.Errorf("Difference with expectation for case %d: %+v", i, merged)
		}
	}
}
