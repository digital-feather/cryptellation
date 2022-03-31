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
		In      []Interface
		InTime  time.Time
		Out     []Interface
		OutTime time.Time
	}{
		{
			In:      []Interface{},
			InTime:  time.Unix(1<<62, 0),
			Out:     []Interface{},
			OutTime: time.Unix(1<<62, 0),
		},
		{
			In: []Interface{
				NewTickEvent(time.Unix(120, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(240, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(180, 0), tick.Tick{}),
			},
			InTime: time.Unix(1<<62, 0),
			Out: []Interface{
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
			},
			OutTime: time.Unix(60, 0),
		},
	}

	for i, c := range cases {
		t, out := OnlyKeepEarliestSameTimeEvents(c.In, c.InTime)
		suite.Require().WithinDuration(c.OutTime, t, time.Microsecond, i)
		suite.Require().Len(out, len(c.Out), i)
	}
}
