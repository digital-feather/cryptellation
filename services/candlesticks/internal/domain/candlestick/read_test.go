package candlestick

import (
	"testing"
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	suite.Suite
}

func (suite *CandlesticksSuite) TestAreMissing() {
	// Given all candlesticks
	cl := NewList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().False(res)
}

func (suite *CandlesticksSuite) TestAreMissingWithOneMissing() {
	// Given all candlesticks
	cl := NewList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().True(res)
}

func (suite *CandlesticksSuite) TestAreMissingWithOneMissingAndLimit() {
	// Given all candlesticks
	cl := NewList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
	}

	// When asking if there is missing candlesticks
	res := AreMissing(cl, time.Unix(0, 0), time.Unix(540, 0), 2)

	// Then there is no missing
	suite.Require().False(res)
}

func (suite *CandlesticksSuite) TestProcessRequestedStartEndTimeWithEmpty() {
	// Given no time

	// When we request the processing
	start, end := ProcessRequestedStartEndTimes(period.M1, nil, nil)

	// Then we have correct times
	suite.Require().WithinDuration(time.Now().Add(-period.M1.Duration()*500), start, time.Minute)
	suite.Require().WithinDuration(time.Now(), end, time.Minute)
}

func (suite *CandlesticksSuite) TestProcessRequestedStartEndTimeWithNoStart() {
	// Given only end
	reqEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := ProcessRequestedStartEndTimes(period.M1, nil, &reqEnd)

	// Then we have correct times
	expectedEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	expectedStart := expectedEnd.Add(-period.M1.Duration() * 500)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}

func (suite *CandlesticksSuite) TestProcessRequestedStartEndTimeWithNoEnd() {
	// Given only start
	reqStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := ProcessRequestedStartEndTimes(period.M1, &reqStart, nil)

	// Then we have correct times
	expectedStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:00Z")
	suite.Require().NoError(err)
	expectedEnd := expectedStart.Add(period.M1.Duration() * 500)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}

func (suite *CandlesticksSuite) TestProcessRequestedStartEndTimeWithNonAligned() {
	// Given a time, but not aligned
	reqStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:04Z")
	suite.Require().NoError(err)
	reqEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:01Z")
	suite.Require().NoError(err)

	// When we request the processing
	start, end := ProcessRequestedStartEndTimes(period.M1, &reqStart, &reqEnd)

	// Then we have correct times
	expectedStart, err := time.Parse(time.RFC3339, "1993-11-15T11:00:00Z")
	suite.Require().NoError(err)
	expectedEnd, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().WithinDuration(expectedStart, start, time.Millisecond)
	suite.Require().WithinDuration(expectedEnd, end, time.Millisecond)
}
