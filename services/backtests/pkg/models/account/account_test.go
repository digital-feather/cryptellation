package account

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

type AccountSuite struct {
	suite.Suite
}

func (suite *AccountSuite) TestValidate() {
	acc := Account{
		Balances: map[string]float64{
			"USDC": 1000,
			"ETH":  1,
		},
	}

	suite.Require().NoError(acc.Validate())
}

func (suite *AccountSuite) TestValidateWithInvalidName() {
	acc := Account{
		Balances: map[string]float64{
			"":    1000,
			"ETH": 1,
		},
	}

	suite.Require().Error(acc.Validate())
}

func (suite *AccountSuite) TestValidateWithInvalidValue() {
	acc := Account{
		Balances: map[string]float64{
			"USDC": 1000,
			"ETH":  -1,
		},
	}

	suite.Require().Error(acc.Validate())
}
