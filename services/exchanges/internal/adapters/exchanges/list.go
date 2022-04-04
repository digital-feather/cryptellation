package exchanges

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
)

var (
	Binance = exchange.Exchange{
		Name: "binance",
		PeriodsSymbols: []string{
			"M1", "M3", "M5", "M15", "M30",
			"H1", "H2", "H4", "H6", "H8", "H12",
			"D1", "D3",
			"W1",
		},
		Fees: 0.1,
	}
)

var (
	Exchanges = []exchange.Exchange{
		Binance,
	}
)
