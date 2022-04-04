package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/account"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
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

func (suite *RedisVdbSuite) TestCreateRead() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	suite.NoError(suite.db.CreateBacktest(context.TODO(), &bt))
	rp, err := suite.db.ReadBacktest(context.TODO(), bt.ID)
	suite.Assert().NoError(err)

	suite.Assert().Equal(bt.ID, rp.ID)
	suite.Assert().Len(rp.Accounts, 1)
	suite.Assert().Len(rp.Accounts["exchange"].Balances, 1)
	suite.Assert().Equal(bt.Accounts["exchange"].Balances["DAI"], rp.Accounts["exchange"].Balances["DAI"])
}

func (suite *RedisVdbSuite) TestUpdate() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.db.CreateBacktest(context.TODO(), &bt))
	bt2 := backtest.Backtest{
		ID: bt.ID,
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	// Should be changes here
	suite.NoError(suite.db.UpdateBacktest(context.TODO(), bt2))
	rp, err := suite.db.ReadBacktest(context.TODO(), bt.ID)
	suite.Assert().NoError(err)

	suite.Equal(bt.ID, rp.ID)
	suite.Equal(bt2.ID, rp.ID)
	suite.Assert().Len(rp.Accounts, 1)
	suite.Assert().Len(rp.Accounts["exchange2"].Balances, 1)
	suite.Assert().Equal(bt2.Accounts["exchange2"].Balances["USDC"], rp.Accounts["exchange2"].Balances["USDC"])
}

func (suite *RedisVdbSuite) TestDelete() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.db.CreateBacktest(context.TODO(), &bt))
	suite.NoError(suite.db.DeleteBacktest(context.TODO(), bt))
	_, err := suite.db.ReadBacktest(context.TODO(), bt.ID)
	suite.Error(err)
}

func (suite *RedisVdbSuite) TestDeleteInexistant() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.db.CreateBacktest(context.TODO(), &bt))
	suite.NoError(suite.db.DeleteBacktest(context.TODO(), bt))
	suite.NoError(suite.db.DeleteBacktest(context.TODO(), bt))
}

func (suite *RedisVdbSuite) TestLock() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.db.CreateBacktest(context.TODO(), &bt))

	// Check that lock works even with panic
	suite.Require().NoError(suite.db.LockedBacktest(bt.ID, func() error {
		panic(errors.New("PANIC"))
	}))

	for i := 0; i < 10; i++ {
		suite.Require().NoError(suite.db.LockedBacktest(bt.ID, func() error {
			return nil
		}), fmt.Sprintf("Lock/Unlock attempt #%d", i+1))
	}

	go suite.db.LockedBacktest(bt.ID, func() error {
		bt.Accounts["exchange"].Balances["ETH"] = 2000
		time.Sleep(200 * time.Millisecond)
		suite.Require().NoError(suite.db.UpdateBacktest(context.TODO(), bt))
		return nil
	})
	time.Sleep(time.Millisecond)

	suite.Require().NoError(suite.db.LockedBacktest(bt.ID, func() error {
		recvBt, err := suite.db.ReadBacktest(context.TODO(), bt.ID)
		suite.Require().NoError(err)
		suite.Require().Equal(2000.0, recvBt.Accounts["exchange"].Balances["ETH"])
		return nil
	}))
}
