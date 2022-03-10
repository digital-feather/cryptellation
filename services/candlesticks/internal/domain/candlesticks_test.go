package domain

import (
	"testing"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	suite.Suite
}

func (suite *CandlesticksSuite) TestAreCsMissing() {
	// Given all candlesticks
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		cl.Set(time.Unix(60*i, 0), candlestick.Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreCsMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().False(res)
}

func (suite *CandlesticksSuite) TestAreCsMissingWithOneMissing() {
	// Given all candlesticks
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		cl.Set(time.Unix(60*i, 0), candlestick.Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreCsMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().True(res)
}

func (suite *CandlesticksSuite) TestAreCsMissingWithOneMissingAndLimit() {
	// Given all candlesticks
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		cl.Set(time.Unix(60*i, 0), candlestick.Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreCsMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 2)

	// Then there is no missing
	suite.Require().False(res)
}
