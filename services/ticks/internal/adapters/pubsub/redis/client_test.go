package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
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

	pairSymbol := "symbol1"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	ch, err := suite.client.Subscribe(ctx, pairSymbol)
	as.NoError(err)

	as.NoError(suite.client.Publish(ctx, t))
	select {
	case recvTick := <-ch:
		as.Equal(t, recvTick)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}
}

func (suite *RedisPubSubSuite) TestOnePubTwoSub() {
	as := suite.Require()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pairSymbol := "symbol2"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ch1, err := suite.client.Subscribe(ctx, pairSymbol)
	as.NoError(err)

	ch2, err := suite.client.Subscribe(ctx, pairSymbol)
	as.NoError(err)

	as.NoError(suite.client.Publish(ctx, t))

	for i := 0; i < 2; i++ {
		select {
		case recvTick := <-ch1:
			as.Equal(t, recvTick)
		case recvTick := <-ch2:
			as.Equal(t, recvTick)
		case <-time.After(1 * time.Second):
			as.FailNow("Timeout")
		}
	}
}

func (suite *RedisPubSubSuite) TestCheckClose() {
	as := suite.Require()

	pairSymbol := "symbol3"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch, err := suite.client.Subscribe(ctx, pairSymbol)
	as.NoError(err)

	cancel()
	as.NoError(suite.client.Publish(context.Background(), t))

	_, open := <-ch
	suite.False(open)
}
