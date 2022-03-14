package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderSideSuite(t *testing.T) {
	suite.Run(t, new(OrderSideTestSuite))
}

type OrderSideTestSuite struct {
	suite.Suite
}

func (suite OrderSideTestSuite) TestValidate() {
	correct1 := Side("buy")
	suite.Assert().NoError(correct1.Validate())
	correct2 := Side("sell")
	suite.Assert().NoError(correct2.Validate())
	incorrect := Side("unknown")
	suite.Assert().Error(incorrect.Validate())
}
