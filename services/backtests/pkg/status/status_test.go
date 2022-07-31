package status

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}

type StatusSuite struct {
	suite.Suite
}

func (suite *StatusSuite) TestMarshalingJSON() {
	as := suite.Require()

	s := Status{
		Finished: true,
	}

	b, err := json.Marshal(s)
	as.NoError(err)

	s2 := Status{}
	as.NoError(json.Unmarshal(b, &s2))
	as.Equal(s, s2)
}
