package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/tick"
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

func (suite *RedisPubSubSuite) BeforeTest(suiteName, testName string) {
	client, err := New()
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *RedisPubSubSuite) TestOnePubOneSubObject() {
	as := suite.Require()
	backtestID := uint(1)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}
	sub, err := suite.client.Subscribe(context.TODO(), backtestID)
	as.NoError(err)
	ch := sub.Channel()

	as.NoError(suite.client.Publish(context.TODO(), backtestID, event.NewTickEvent(ts, t)))
	select {
	case recvEvent := <-ch:
		suite.checkTick(recvEvent, ts, t)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}

	as.NoError(suite.client.Publish(context.TODO(), backtestID, event.NewEndEvent(ts)))
	select {
	case recvEvent := <-ch:
		suite.checkEnd(recvEvent, ts)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}

	as.NoError(sub.Close())
}

func (suite *RedisPubSubSuite) TestOnePubTwoSub() {
	as := suite.Require()
	backtestID := uint(1)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}

	sub1, err := suite.client.Subscribe(context.TODO(), backtestID)
	as.NoError(err)
	ch1 := sub1.Channel()

	sub2, err := suite.client.Subscribe(context.TODO(), backtestID)
	as.NoError(err)
	ch2 := sub2.Channel()

	as.NoError(suite.client.Publish(context.TODO(), backtestID, event.NewTickEvent(ts, t)))

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

	as.NoError(sub1.Close())
	as.NoError(sub2.Close())
}

func (suite *RedisPubSubSuite) TestCheckClose() {
	as := suite.Require()
	backtestID := uint(1)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}

	sub, err := suite.client.Subscribe(context.TODO(), backtestID)
	as.NoError(err)
	ch := sub.Channel()

	as.NoError(sub.Close())

	_, open := <-ch
	suite.False(open)

	ch2 := sub.Channel()
	_, open = <-ch2
	suite.False(open)

	as.Error(sub.Close())

	as.NoError(suite.client.Publish(context.TODO(), backtestID, event.NewTickEvent(ts, t)))
}

func (suite *RedisPubSubSuite) checkTick(evt event.Interface, t time.Time, ti tick.Tick) {
	as := suite.Require()

	as.Equal(event.TypeIsTick, evt.GetType())
	rt, ok := evt.(event.TickEvent)
	as.True(ok)
	as.Equal(ti, rt.Content)
}

func (suite *RedisPubSubSuite) checkEnd(evt event.Interface, t time.Time) {
	as := suite.Require()

	as.Equal(event.TypeIsEnd, evt.GetType())
	as.Equal(t, evt.GetTime())
	_, ok := evt.(event.EndEvent)
	as.True(ok)
}
