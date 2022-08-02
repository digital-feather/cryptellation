package cockroach

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

func TestCockroachDatabaseSuite(t *testing.T) {
	if os.Getenv("COCKROACHDB_HOST") == "" {
		t.Skip()
	}

	suite.Run(t, new(CockroachDatabaseSuite))
}

type CockroachDatabaseSuite struct {
	suite.Suite
	db *DB
}

func (suite *CockroachDatabaseSuite) SetupTest() {
	defer tmpEnvVar("COCKROACHDB_DATABASE", "candlesticks")()

	db, err := New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.db = db
}

func (suite *CockroachDatabaseSuite) TestNewWithURIError() {
	defer tmpEnvVar("COCKROACHDB_HOST", "")()

	var err error
	_, err = New()
	suite.Require().Error(err)
}

func (suite *CockroachDatabaseSuite) TestCreate() {
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t, candlestick.Candlestick{
		Open:       1,
		Low:        0.5,
		High:       2,
		Close:      1.5,
		Volume:     1000,
		Uncomplete: true,
	}))
	recvList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))
	err = suite.db.ReadCandlesticks(context.Background(), recvList, t.Add(-time.Hour), t.Add(time.Hour), 0)
	suite.Require().NoError(err)

	suite.Require().Equal(1, recvList.Len())
	cs, _ := list.Get(t)
	rcs, exists := recvList.Get(t)
	suite.True(exists)
	suite.True(cs == rcs)
}

func (suite *CockroachDatabaseSuite) TestCreateTwice() {
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t, candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))
	suite.Require().Error(suite.db.CreateCandlesticks(context.Background(), list))
}

func (suite *CockroachDatabaseSuite) TestRead() {
	// Create targeted exchange, pair, period
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	t1, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t1, candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	t2, err := time.Parse(time.RFC3339, "1993-11-15T11:30:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t2, candlestick.Candlestick{
		Open:   2,
		Low:    1,
		High:   4,
		Close:  3,
		Volume: 2000,
	}))

	t3, err := time.Parse(time.RFC3339, "1993-11-15T11:31:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t3, candlestick.Candlestick{
		Open:   3,
		Low:    1.5,
		High:   6,
		Close:  4.5,
		Volume: 3000,
	}))

	t4, err := time.Parse(time.RFC3339, "1993-11-15T11:32:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t4, candlestick.Candlestick{
		Open:   4,
		Low:    2,
		High:   8,
		Close:  6,
		Volume: 4000,
	}))

	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))

	// Create other exchange
	otherExchangeList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange2",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(err)
	suite.Require().NoError(otherExchangeList.Set(t2, candlestick.Candlestick{
		Open:   1,
		Low:    1,
		High:   1,
		Close:  1,
		Volume: 1,
	}))
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), otherExchangeList))

	// Create other pair
	otherPairList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "BTC-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(err)
	suite.Require().NoError(otherPairList.Set(t2, candlestick.Candlestick{
		Open:   2,
		Low:    2,
		High:   2,
		Close:  2,
		Volume: 2,
	}))
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), otherPairList))

	// Create other period
	otherPeriodList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M15,
	})
	suite.Require().NoError(err)
	suite.Require().NoError(otherPeriodList.Set(t2, candlestick.Candlestick{
		Open:   3,
		Low:    3,
		High:   3,
		Close:  3,
		Volume: 3,
	}))
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), otherPeriodList))

	// Read only the two centered candles
	recvList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t2, t3, 0))
	suite.Require().Equal(2, recvList.Len())
	c, _ := list.Get(t2)
	rc, exists := recvList.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	c, _ = list.Get(t3)
	rc, exists = recvList.Get(t3)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	// Check others
	recvList = candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange2",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t2, t3, 0))
	c, _ = otherExchangeList.Get(t2)
	rc, exists = recvList.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	recvList = candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "BTC-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t2, t3, 0))
	c, _ = otherPairList.Get(t2)
	rc, exists = recvList.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	recvList = candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M15,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t2, t3, 0))
	c, _ = otherPeriodList.Get(t2)
	rc, exists = recvList.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)
}

func (suite *CockroachDatabaseSuite) TestReadLimit() {
	// Create targeted exchange, pair, period
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})

	t1, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t1, candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	t2, err := time.Parse(time.RFC3339, "1993-11-15T11:30:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t2, candlestick.Candlestick{
		Open:   2,
		Low:    1,
		High:   4,
		Close:  3,
		Volume: 2000,
	}))

	t3, err := time.Parse(time.RFC3339, "1993-11-15T11:31:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t3, candlestick.Candlestick{
		Open:   3,
		Low:    1.5,
		High:   6,
		Close:  4.5,
		Volume: 3000,
	}))

	t4, err := time.Parse(time.RFC3339, "1993-11-15T11:32:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t4, candlestick.Candlestick{
		Open:   4,
		Low:    2,
		High:   8,
		Close:  6,
		Volume: 4000,
	}))

	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))

	// Read only the 2 first
	recvList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t1, t4, 2))
	suite.Require().Equal(2, recvList.Len())
	c, _ := list.Get(t1)
	rc, exists := recvList.Get(t1)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	c, _ = list.Get(t2)
	rc, exists = recvList.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)
}

func (suite *CockroachDatabaseSuite) TestReadEmpty() {
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	recvList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(suite.db.ReadCandlesticks(context.Background(), recvList, t.Add(-time.Hour), t.Add(time.Hour), 0))
	suite.Require().Equal(0, recvList.Len())
}

func (suite *CockroachDatabaseSuite) TestUpdate() {
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t, candlestick.Candlestick{
		Open:       1,
		Low:        0.5,
		High:       2,
		Close:      1.5,
		Volume:     1000,
		Uncomplete: true,
	}))
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))

	update := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(update.Set(t, candlestick.Candlestick{
		Open:       2,
		Low:        1,
		High:       4,
		Close:      3,
		Volume:     2000,
		Uncomplete: false,
	}))
	suite.Require().NoError(suite.db.UpdateCandlesticks(context.Background(), update))
	receivedList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	err = suite.db.ReadCandlesticks(context.Background(), receivedList, t.Add(-time.Hour), t.Add(time.Hour), 0)
	suite.Require().NoError(err)

	suite.Require().Equal(1, receivedList.Len())
	cs, _ := update.Get(t)
	rcs, exists := receivedList.Get(t)
	suite.True(exists)
	suite.Require().Equal(cs, rcs)
}

func (suite *CockroachDatabaseSuite) TestUpdateInexistantTwice() {
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	suite.Require().NoError(list.Set(t, candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	suite.Require().Error(suite.db.UpdateCandlesticks(context.Background(), list))
	suite.Require().Error(suite.db.UpdateCandlesticks(context.Background(), list))
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	list := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	for i := 0; i < 10; i++ {
		suite.Require().NoError(list.Set(time.Unix(int64(i*int(period.M1.Duration().Seconds())), 0), candlestick.Candlestick{
			Open:   1 + float64(i),
			Low:    0.5 + float64(i),
			High:   2 + float64(i),
			Close:  1.5 + float64(i),
			Volume: 1000 * float64(i),
		}))
	}
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), list))

	// Remove half the data
	delete := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(delete.Set(time.Unix(int64(0*period.M1.Duration().Seconds()), 0), candlestick.Candlestick{}))
	suite.Require().NoError(delete.Set(time.Unix(int64(1*period.M1.Duration().Seconds()), 0), candlestick.Candlestick{}))
	suite.Require().NoError(delete.Set(time.Unix(int64(2*period.M1.Duration().Seconds()), 0), candlestick.Candlestick{}))
	suite.Require().NoError(delete.Set(time.Unix(int64(3*period.M1.Duration().Seconds()), 0), candlestick.Candlestick{}))
	suite.Require().NoError(delete.Set(time.Unix(int64(4*period.M1.Duration().Seconds()), 0), candlestick.Candlestick{}))
	suite.Require().NoError(suite.db.DeleteCandlesticks(context.Background(), delete))

	// Check staying data
	tEnd := time.Unix(10*int64(period.M1.Duration().Seconds()), 0)
	tStart := time.Unix(5*int64(period.M1.Duration().Seconds()), 0)
	receivedList := candlestick.NewList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	err := suite.db.ReadCandlesticks(context.Background(), receivedList, tStart, tEnd, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(list.Len()-5, receivedList.Len())
	suite.Require().NoError(receivedList.Loop(func(ts time.Time, cs candlestick.Candlestick) (bool, error) {
		origCS, exists := list.Get(ts)
		suite.True(exists)
		suite.Require().Equal(origCS, cs)
		return false, nil
	}))
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
