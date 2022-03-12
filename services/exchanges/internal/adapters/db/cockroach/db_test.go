package cockroach

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/services/exchanges/pkg/exchange"
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

func (suite *CockroachDatabaseSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "exchanges")()

	db, err := New()
	suite.Require().NoError(err)
	suite.db = db

	suite.Require().NoError(Reset())
}

func (suite *CockroachDatabaseSuite) TestNewWithURIError() {
	defer tests.TempEnvVar("COCKROACHDB_HOST", "")()

	var err error
	_, err = New()
	suite.Error(err)
}

func (suite *CockroachDatabaseSuite) TestCreateRead() {
	as := suite.Require()

	// Given a exchange
	p := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC(),
	}

	// When we create it and read it
	as.NoError(suite.db.CreateExchanges(context.Background(), p))
	rp, err := suite.db.ReadExchanges(context.Background(), p.Name)
	as.NoError(err)

	// Then it's the same
	as.Len(rp, 1)
	as.Equal(p.Name, rp[0].Name)
	as.Contains(rp[0].Pairs, "ABC-DEF")
	as.Contains(rp[0].Pairs, "IJK-LMN")
	as.Contains(rp[0].Periods, "M1")
	as.Contains(rp[0].Periods, "M3")
	as.Equal(p.Fees, rp[0].Fees)
	as.WithinDuration(time.Now().UTC(), rp[0].LastSyncTime, time.Second)
}

func (suite *CockroachDatabaseSuite) TestCreateReadInexistant() {
	as := suite.Require()

	// When we read an inexistant exchange
	exchanges, err := suite.db.ReadExchanges(context.Background(), "inexistant")

	// Then there is no error but no exchange
	as.NoError(err)
	as.Len(exchanges, 0)
}

func (suite *CockroachDatabaseSuite) TestReadAll() {
	as := suite.Require()

	// Given 3 created exchanges
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p1))
	p2 := exchange.Exchange{
		Name:         "exchange2",
		Pairs:        []string{"ABC-DEF", "XYZ-LMN"},
		Periods:      []string{"M1", "M5"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p2))
	p3 := exchange.Exchange{
		Name:         "exchange3",
		Pairs:        []string{"ABC-XYZ", "IJK-LMN"},
		Periods:      []string{"M1", "M10"},
		Fees:         0.3,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p3))

	// When we read all of them
	ps, err := suite.db.ReadExchanges(context.Background())
	as.NoError(err)

	// Then we have all of them
	as.Len(ps, 3)
	for _, p := range ps {
		if p.Name != p1.Name && p.Name != p2.Name && p.Name != p3.Name {
			as.Fail("This exchange should not exists", p)
		}
	}
}

func (suite *CockroachDatabaseSuite) TestUpdate() {
	as := suite.Require()

	// Given a created exchange
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC().Add(-time.Hour),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p1))

	// When we make some modification to it
	p2 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}

	// And we update it
	as.NoError(suite.db.UpdateExchanges(context.Background(), p2))

	// Then the exchange is updated
	rp, err := suite.db.ReadExchanges(context.Background(), p2.Name)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(p2.Name, rp[0].Name)
	as.Contains(rp[0].Pairs, "ABC-XYZ")
	as.Contains(rp[0].Pairs, "IJK-XYZ")
	as.Contains(rp[0].Periods, "M5")
	as.Contains(rp[0].Periods, "D1")
	as.Equal(p2.Fees, rp[0].Fees)
	as.WithinDuration(time.Now().UTC(), rp[0].LastSyncTime, time.Second)
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	as := suite.Require()

	// Given twp created exchange
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p1))
	p2 := exchange.Exchange{
		Name:         "exchange2",
		Pairs:        []string{"ABC-XYZ", "IJK-ABC"},
		Periods:      []string{"M5", "M1"},
		Fees:         0.3,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p2))

	// When we delete it
	as.NoError(suite.db.DeleteExchanges(context.Background(), p1.Name))

	// Then we can't read it anymore
	exchanges, err := suite.db.ReadExchanges(context.Background(), p1.Name, p2.Name)
	as.NoError(err)
	as.Len(exchanges, 1)
	as.Equal(p2.Name, exchanges[0].Name)

	// And there is no pair left
	pairs := []Pair{}
	as.NoError(suite.db.client.WithContext(context.Background()).Find(&pairs).Error)
	as.Len(pairs, 2)
	as.Contains(pairs, Pair{Symbol: "ABC-XYZ"})
	as.Contains(pairs, Pair{Symbol: "IJK-ABC"})

	// And there is no period left
	periods := []Period{}
	as.NoError(suite.db.client.WithContext(context.Background()).Find(&periods).Error)
	as.Len(periods, 2)
	as.Contains(periods, Period{Symbol: "M5"})
	as.Contains(periods, Period{Symbol: "M1"})
}

func (suite *CockroachDatabaseSuite) TestReset() {
	as := suite.Require()

	// Given a created exchange
	p := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.db.CreateExchanges(context.Background(), p))

	// When we reset the DB
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "exchanges")()
	as.NoError(Reset())

	// Then there is no exchange left
	exchanges := []Exchange{}
	as.NoError(suite.db.client.WithContext(context.Background()).Find(&exchanges).Error)
	as.Len(exchanges, 0)

	// And there is no pair left
	pairs := []Pair{}
	as.NoError(suite.db.client.WithContext(context.Background()).Find(&pairs).Error)
	as.Len(pairs, 0)

	// And there is no period left
	periods := []Period{}
	as.NoError(suite.db.client.WithContext(context.Background()).Find(&periods).Error)
	as.Len(periods, 0)
}
