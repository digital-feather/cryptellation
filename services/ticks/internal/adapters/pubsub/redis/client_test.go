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

func (suite *RedisPubSubSuite) BeforeTest(suiteName, testName string) {
	client, err := New()
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *RedisPubSubSuite) TestOnePubOneSubObject() {
	as := suite.Require()
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}
	sub, err := suite.client.Subscribe(context.TODO(), "BTC-USDC")
	as.NoError(err)
	ch := sub.Channel()

	as.NoError(suite.client.Publish(context.TODO(), t))
	select {
	case recvTick := <-ch:
		as.Equal(t, recvTick)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}

	as.NoError(sub.Close())
}

func (suite *RedisPubSubSuite) TestOnePubTwoSub() {
	as := suite.Require()
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}

	sub1, err := suite.client.Subscribe(context.TODO(), "BTC-USDC")
	as.NoError(err)
	ch1 := sub1.Channel()

	sub2, err := suite.client.Subscribe(context.TODO(), "BTC-USDC")
	as.NoError(err)
	ch2 := sub2.Channel()

	as.NoError(suite.client.Publish(context.TODO(), t))

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

	as.NoError(sub1.Close())
	as.NoError(sub2.Close())
}

func (suite *RedisPubSubSuite) TestCheckClose() {
	as := suite.Require()
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: "BTC-USDC",
		Price:      40000,
		Exchange:   "exchange",
	}

	sub, err := suite.client.Subscribe(context.TODO(), "BTC-USDC")
	as.NoError(err)
	ch := sub.Channel()

	as.NoError(sub.Close())

	_, open := <-ch
	suite.False(open)

	ch2 := sub.Channel()
	_, open = <-ch2
	suite.False(open)

	as.Error(sub.Close())

	as.NoError(suite.client.Publish(context.TODO(), t))
}
