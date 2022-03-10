package binance

import (
	"testing"

	"github.com/cryptellation/cryptellation/pkg/types/period"
)

var possibleIntervals = map[period.Symbol]string{
	period.M1:  "1m",
	period.M3:  "3m",
	period.M5:  "5m",
	period.M15: "15m",
	period.M30: "30m",
	period.H1:  "1h",
	period.H2:  "2h",
	period.H4:  "4h",
	period.H6:  "6h",
	period.H8:  "8h",
	period.H12: "12h",
	period.D1:  "1d",
	period.D3:  "3d",
	period.W1:  "1w",
}

func TestPeriodToInterval(t *testing.T) {
	for k, v := range possibleIntervals {
		if e, err := PeriodToInterval(k); err != nil {
			t.Error("Interval for Period", k, "should not throw an error:", err)
		} else if e != v {
			t.Error("Interval for Period", k, "does not correspond : should be", v, "but is", e)
		}
	}
}

func TestPeriodToInterval_InexistantPeriod(t *testing.T) {
	if _, err := PeriodToInterval("unknown"); err == nil {
		t.Error("Period 0 should throw an error")
	}
}
