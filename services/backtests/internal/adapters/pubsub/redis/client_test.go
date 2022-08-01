package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

func TestRedisPubSubSuite(t *testing.T) {
	if os.Getenv("REDIS_ADDRESS") == "" {
		t.Skip()
	}

	suite.Run(t, new(RedisPubSubSuite))
}

type RedisPubSubSuite struct {
	suite.Suite
	client *Client
}

func (suite *RedisPubSubSuite) SetupTest() {
	client, err := New()
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *RedisPubSubSuite) TestOnePubOneSubObject() {
	as := suite.Require()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backtestID := uint(1)
	ts := time.Unix(60, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	st := status.Status{
		Finished: true,
	}

	ch, err := suite.client.Subscribe(ctx, backtestID)
	as.NoError(err)

	as.NoError(suite.client.Publish(ctx, backtestID, event.NewTickEvent(ts, t)))
	select {
	case recvEvent := <-ch:
		suite.checkTick(recvEvent, ts, t)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}

	as.NoError(suite.client.Publish(ctx, backtestID, event.NewStatusEvent(ts, st)))
	select {
	case recvEvent := <-ch:
		suite.checkEnd(recvEvent, ts, st)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}
}

func (suite *RedisPubSubSuite) TestOnePubTwoSub() {
	as := suite.Require()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backtestID := uint(2)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ch1, err := suite.client.Subscribe(ctx, backtestID)
	as.NoError(err)

	ch2, err := suite.client.Subscribe(ctx, backtestID)
	as.NoError(err)

	as.NoError(suite.client.Publish(ctx, backtestID, event.NewTickEvent(ts, t)))

	for i := 0; i < 2; i++ {
		select {
		case recvEvent := <-ch1:
			suite.checkTick(recvEvent, ts, t)
		case recvEvent := <-ch2:
			suite.checkTick(recvEvent, ts, t)
		case <-time.After(1 * time.Second):
			as.FailNow("Timeout")
		}
	}
}

func (suite *RedisPubSubSuite) TestCheckClose() {
	as := suite.Require()

	backtestID := uint(3)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch, err := suite.client.Subscribe(ctx, backtestID)
	as.NoError(err)

	cancel()
	as.NoError(suite.client.Publish(context.Background(), backtestID, event.NewTickEvent(ts, t)))

	_, open := <-ch
	suite.False(open)
}

func (suite *RedisPubSubSuite) checkTick(evt event.Event, t time.Time, ti tick.Tick) {
	as := suite.Require()

	as.Equal(event.TypeIsTick, evt.Type)
	as.Equal(t, evt.Time)
	rt, ok := evt.Content.(tick.Tick)
	as.True(ok)
	as.Equal(ti, rt)
}

func (suite *RedisPubSubSuite) checkEnd(evt event.Event, t time.Time, st status.Status) {
	as := suite.Require()

	as.Equal(event.TypeIsStatus, evt.Type)
	as.Equal(t, evt.Time)
	rt, ok := evt.Content.(status.Status)
	as.True(ok)
	as.Equal(st, rt)
}
