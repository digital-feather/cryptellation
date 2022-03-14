package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderTypeSuite(t *testing.T) {
	suite.Run(t, new(OrderTypeTestSuite))
}

type OrderTypeTestSuite struct {
	suite.Suite
}

func (suite OrderTypeTestSuite) TestValidate() {
	correct := Type("market")
	suite.Assert().NoError(correct.Validate())
	incorrect := Type("unknown")
	suite.Assert().Error(incorrect.Validate())
}
