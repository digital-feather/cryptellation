package tick

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTickSuite(t *testing.T) {
	suite.Run(t, new(TickSuite))
}

type TickSuite struct {
	suite.Suite
}

func (suite *TickSuite) TestMarshalingJSON() {
	as := suite.Require()

	tick := Tick{
		PairSymbol: "BTC-USDC",
		Price:      1.01,
		Exchange:   "exchange",
	}

	b, err := json.Marshal(tick)
	as.NoError(err)

	tick2 := Tick{}
	as.NoError(json.Unmarshal(b, &tick2))
	as.Equal(tick, tick2)
}
