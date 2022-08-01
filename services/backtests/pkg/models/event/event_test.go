package event

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
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
		In      []Event
		InTime  time.Time
		Out     []Event
		OutTime time.Time
	}{
		{
			In:      []Event{},
			InTime:  time.Unix(1<<62, 0),
			Out:     []Event{},
			OutTime: time.Unix(1<<62, 0),
		},
		{
			In: []Event{
				NewTickEvent(time.Unix(120, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(240, 0), tick.Tick{}),
				NewTickEvent(time.Unix(60, 0), tick.Tick{}),
				NewTickEvent(time.Unix(180, 0), tick.Tick{}),
			},
			InTime: time.Unix(1<<62, 0),
			Out: []Event{
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
