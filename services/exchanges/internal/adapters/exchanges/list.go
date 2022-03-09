package exchanges

import (
	"github.com/cryptellation/cryptellation/pkg/types/exchange"
	"github.com/cryptellation/cryptellation/pkg/types/period"
)

var (
	Binance = exchange.Exchange{
		Name: "binance",
		Periods: []string{
			period.M1.String(),
			period.M3.String(),
			period.M5.String(),
			period.M15.String(),
			period.M30.String(),
			period.H1.String(),
			period.H2.String(),
			period.H4.String(),
			period.H6.String(),
			period.H8.String(),
			period.H12.String(),
			period.D1.String(),
			period.D3.String(),
			period.W1.String(),
		},
		Fees: 0.1,
	}
)

var (
	Exchanges = []exchange.Exchange{
		Binance,
	}
)
