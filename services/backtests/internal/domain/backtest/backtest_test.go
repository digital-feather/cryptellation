package backtest

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/account"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	suite.Suite
}

func (suite *BacktestSuite) TestMarshalUnMarshalBinary() {
	bt := Backtest{
		StartTime: time.Unix(0, 0).UTC(),
		CurrentCsTick: CurrentCsTick{
			Time: time.Unix(100, 0).UTC(),
		},
		EndTime: time.Unix(200, 0).UTC(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"USDC": 10000,
				},
			},
		},
	}

	content, err := bt.MarshalBinary()
	suite.Require().NoError(err)
	nbt := Backtest{}
	suite.Require().NoError(nbt.UnmarshalBinary(content))
	suite.Require().Equal(bt, nbt)
}

func (suite *BacktestSuite) TestIncrementPriceID() {
	// TODO
}
