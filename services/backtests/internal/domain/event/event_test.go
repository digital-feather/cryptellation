package event

import (
	"testing"
	"time"

	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/tick"
	"github.com/stretchr/testify/suite"
)

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

type EventSuite struct {
	suite.Suite
}

func (suite *EventSuite) TestOnlyKeepEarliestSameTimeEvents() {
	cases := []struct {
		In   []Interface
		Time time.Time
		Out  []Interface
	}{
		{
			In:   []Interface{},
			Time: time.Unix(1<<62, 0),
			Out:  []Interface{},
		},
		{
			In: []Interface{
				NewTickEvent(time.Unix(120, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(240, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(180, 0), tick.Tick{}),
			},
			Time: time.Unix(60, 0),
			Out: []Interface{
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
			},
		},
	}

	for i, c := range cases {
		t, out := OnlyKeepEarliestSameTimeEvents(c.In)
		suite.Require().WithinDuration(c.Time, t, time.Microsecond, i)
		suite.Require().Len(out, len(c.Out), i)
	}
}
