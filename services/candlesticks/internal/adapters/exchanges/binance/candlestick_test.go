package binance

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/candlestick"
	"github.com/cryptellation/cryptellation/services/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
)

func TestBinanceSuite(t *testing.T) {
	if os.Getenv("BINANCE_API_KEY") == "" {
		t.Skip()
	}

	suite.Run(t, new(BinanceSuite))
}

type BinanceSuite struct {
	suite.Suite
	service *Service
}

func (suite *BinanceSuite) BeforeTest(suiteName, testName string) {
	service, err := New()
	suite.Require().NoError(err)
	suite.service = service
}

func (suite *BinanceSuite) TestCandlesticks() {
	p := "BTC-USDT"
	s, err := suite.service.Candlesticks(p, period.D1)
	suite.Require().NoError(err)
	suite.Require().NotNil(s)
	suite.Require().Equal(period.D1, s.Period())
	suite.Require().Equal(p, s.PairSymbol())
}

func (suite *BinanceSuite) TestDo() {
	p := "BTC-USDC"
	s, err := suite.service.Candlesticks(p, period.M1)
	suite.Require().NoError(err)

	ts, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:00:00")
	suite.Require().NoError(err)

	te, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:05:00")
	suite.Require().NoError(err)

	cs, err := s.Limit(2).StartTime(ts).EndTime(te).Do(context.TODO())
	suite.Require().NoError(err)

	suite.Require().Equal(p, cs.PairSymbol())
	suite.Require().Equal(period.M1, cs.Period())

	expected := candlestick.Candlestick{
		Open:   16084.16,
		High:   16093.26,
		Low:    16084.16,
		Close:  16093.26,
		Volume: 0.344592,
	}

	suite.Require().Equal(2, cs.Len())
	rc, exists := cs.Get(ts)
	suite.Require().True(exists)
	suite.Require().Equal(expected, rc)
}

func (suite *BinanceSuite) TestDoWithUncompleteCandlestick() {
	// TODO
}
