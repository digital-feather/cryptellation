package cockroach

import (
	"context"
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/services/pairs/internal/domain/pair"
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
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "pairs")()

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

	// Given a pair
	p := pair.Pair{
		BaseAssetSymbol:  "ETH",
		QuoteAssetSymbol: "USDC",
	}

	// When we create it and read it
	as.NoError(suite.db.CreatePairs(context.Background(), p))
	rp, err := suite.db.ReadPairs(context.Background(), p.Symbol())
	as.NoError(err)

	// Then it's the same
	as.Len(rp, 1)
	as.Equal(p, rp[0])
}

func (suite *CockroachDatabaseSuite) TestCreateReadInexistant() {
	as := suite.Require()

	// When we read an inexistant pair
	p := pair.Pair{
		BaseAssetSymbol:  "BTC",
		QuoteAssetSymbol: "USDC",
	}
	pairs, err := suite.db.ReadPairs(context.Background(), p.Symbol())

	// Then there is no error but no pair
	as.NoError(err)
	as.Len(pairs, 0)
}

func (suite *CockroachDatabaseSuite) TestReadAll() {
	as := suite.Require()

	// Given 3 created pairs
	p1 := pair.Pair{
		BaseAssetSymbol:  "ETH",
		QuoteAssetSymbol: "USDC",
	}
	as.NoError(suite.db.CreatePairs(context.Background(), p1))
	p2 := pair.Pair{
		BaseAssetSymbol:  "BTC",
		QuoteAssetSymbol: "USDC",
	}
	as.NoError(suite.db.CreatePairs(context.Background(), p2))
	p3 := pair.Pair{
		BaseAssetSymbol:  "FTM",
		QuoteAssetSymbol: "USDC",
	}
	as.NoError(suite.db.CreatePairs(context.Background(), p3))

	// When we read all of them
	ps, err := suite.db.ReadPairs(context.Background())
	as.NoError(err)

	// Then we have all of them
	as.Len(ps, 3)
	for _, p := range ps {
		if p != p1 && p != p2 && p != p3 {
			as.Fail("This pair should not exists", p)
		}
	}
}

func (suite *CockroachDatabaseSuite) TestUpdate() {
	as := suite.Require()

	// Given a created pair
	p1 := pair.Pair{
		BaseAssetSymbol:  "ETH",
		QuoteAssetSymbol: "USDC",
	}
	as.NoError(suite.db.CreatePairs(context.Background(), p1))

	// When we make some modification to it
	p2 := p1
	// TODO: Should be changes here

	// And we update it
	as.NoError(suite.db.UpdatePairs(context.Background(), p2))

	// Then the pair is updated
	rp, err := suite.db.ReadPairs(context.Background(), p2.Symbol())
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(p2, rp[0])
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	as := suite.Require()

	// Given a created pair
	p := pair.Pair{
		BaseAssetSymbol:  "ETH",
		QuoteAssetSymbol: "USDC",
	}
	as.NoError(suite.db.CreatePairs(context.Background(), p))

	// When we delete it
	as.NoError(suite.db.DeletePairs(context.Background(), p.Symbol()))

	// Then we can't read it anymore
	pairs, err := suite.db.ReadPairs(context.Background(), p.Symbol())
	as.NoError(err)
	as.Len(pairs, 0)
}
