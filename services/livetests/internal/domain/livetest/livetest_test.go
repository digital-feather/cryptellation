package livetest

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/account"
	"github.com/stretchr/testify/suite"
)

func TestLivetestSuite(t *testing.T) {
	suite.Run(t, new(LivetestSuite))
}

type LivetestSuite struct {
	suite.Suite
}

func (suite *LivetestSuite) TestMarshalUnMarshalBinary() {
	bt := Livetest{
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
	nbt := Livetest{}
	suite.Require().NoError(nbt.UnmarshalBinary(content))
	suite.Require().Equal(bt, nbt)
}
