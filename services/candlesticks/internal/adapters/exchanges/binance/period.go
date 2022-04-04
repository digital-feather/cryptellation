package binance

import (
	"fmt"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/period"
)

// Intervals represents every intervals supported by Binance API
func Intervals() []period.Symbol {
	return []period.Symbol{
		period.M1,
		period.M3,
		period.M5,
		period.M15,
		period.M30,
		period.H1,
		period.H2,
		period.H4,
		period.H6,
		period.H8,
		period.H12,
		period.D1,
		period.D3,
		period.W1,
	}
}

// PeriodToInterval converts an interval to its corresponding epoch
func PeriodToInterval(interval period.Symbol) (e string, err error) {
	switch interval {
	case period.M1:
		return "1m", nil
	case period.M3:
		return "3m", nil
	case period.M5:
		return "5m", nil
	case period.M15:
		return "15m", nil
	case period.M30:
		return "30m", nil
	case period.H1:
		return "1h", nil
	case period.H2:
		return "2h", nil
	case period.H4:
		return "4h", nil
	case period.H6:
		return "6h", nil
	case period.H8:
		return "8h", nil
	case period.H12:
		return "12h", nil
	case period.D1:
		return "1d", nil
	case period.D3:
		return "3d", nil
	case period.W1:
		return "1w", nil
	default:
		return e, fmt.Errorf("interval error: unknown period")
	}
}
