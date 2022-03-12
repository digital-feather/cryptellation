package exchange

import (
	"reflect"
	"testing"

	"github.com/cryptellation/cryptellation/services/candlesticks/pkg/period"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		Exchange1 Exchange
		Exchange2 Exchange
		Expected  Exchange
	}{
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{period.M1.String(), period.M15.String()},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange2",
				Pairs:   []string{"ABC-DEF", "ABC-XYZ"},
				Periods: []string{period.M1.String(), period.M3.String()},
				Fees:    0.2,
			},
			Expected: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ", "ABC-XYZ"},
				Periods: []string{period.M1.String(), period.M15.String(), period.M3.String()},
				Fees:    0.1,
			},
		},
		{
			Exchange1: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{period.M1.String(), period.M15.String()},
				Fees:    0.1,
			},
			Exchange2: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{period.M1.String(), period.M15.String()},
				Fees:    0.1,
			},
			Expected: Exchange{
				Name:    "exchange1",
				Pairs:   []string{"ABC-DEF", "DEF-XYZ"},
				Periods: []string{period.M1.String(), period.M15.String()},
				Fees:    0.1,
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
