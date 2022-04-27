package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestRedisVdbSuite(t *testing.T) {
	if os.Getenv("REDIS_ADDRESS") == "" {
		t.Skip()
	}

	suite.Run(t, new(RedisVdbSuite))
}

type RedisVdbSuite struct {
	suite.Suite
	db *DB
}

func (suite *RedisVdbSuite) BeforeTest(suiteName, testName string) {
	db, err := New()
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *RedisVdbSuite) ResetSymbolListener(exchange, pairSymbol string) {
	suite.db.client.Set(context.Background(), fmt.Sprintf(redisKeySymbolListener, exchange, pairSymbol), 0, time.Second)
}

func (suite *RedisVdbSuite) TestIncrementDecrementSymbolListener() {
	suite.ResetSymbolListener("EXCHANGE1", "PAIR1")
	suite.ResetSymbolListener("EXCHANGE2", "PAIR1")

	count, err := suite.db.IncrementSymbolListenerCount(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.db.IncrementSymbolListenerCount(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(2), count)

	count, err = suite.db.IncrementSymbolListenerCount(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.db.GetSymbolListenerCount(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(2), count)

	count, err = suite.db.GetSymbolListenerCount(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.db.DecrementSymbolListenerCount(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	suite.NoError(suite.db.ClearSymbolListenersCount(context.Background()))

	count, err = suite.db.GetSymbolListenerCount(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(0), count)

	count, err = suite.db.GetSymbolListenerCount(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(0), count)
}
