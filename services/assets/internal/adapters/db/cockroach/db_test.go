package cockroach

import (
	"context"
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
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
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "assets")()

	db, _, err := New()
	suite.Require().NoError(err)
	suite.db = db

	suite.Require().NoError(Reset())
}

func (suite *CockroachDatabaseSuite) AfterTest(suiteName, testName string) {
	suite.db.Close()
}

func (suite *CockroachDatabaseSuite) TestNewWithURIError() {
	defer tests.TempEnvVar("COCKROACHDB_HOST", "")()

	var err error
	_, _, err = New()
	suite.Error(err)
}

func (suite *CockroachDatabaseSuite) TestCreateRead() {
	as := suite.Require()

	// Given an asset
	a := asset.Asset{Symbol: "ETH"}

	// When we create it and read it
	as.NoError(suite.db.CreateAssets(context.Background(), a))
	rp, err := suite.db.ReadAssets(context.Background(), a.Symbol)
	as.NoError(err)

	// Then it's the same
	as.Len(rp, 1)
	as.Equal(a, rp[0])
}

func (suite *CockroachDatabaseSuite) TestCreateReadInexistant() {
	as := suite.Require()

	// When we read an inexistant asset
	a := asset.Asset{Symbol: "BTC"}
	assets, err := suite.db.ReadAssets(context.Background(), a.Symbol)

	// Then there is no error but no asset
	as.NoError(err)
	as.Len(assets, 0)
}

func (suite *CockroachDatabaseSuite) TestReadAll() {
	as := suite.Require()

	// Given 3 created assets
	a1 := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a1))
	a2 := asset.Asset{Symbol: "FTM"}
	as.NoError(suite.db.CreateAssets(context.Background(), a2))
	a3 := asset.Asset{Symbol: "DAI"}
	as.NoError(suite.db.CreateAssets(context.Background(), a3))

	// When we read all of them
	ps, err := suite.db.ReadAssets(context.Background())
	as.NoError(err)

	// Then we have all of them
	as.Len(ps, 3)
	for _, p := range ps {
		if p.Symbol != a1.Symbol && p.Symbol != a2.Symbol && p.Symbol != a3.Symbol {
			as.Fail("This asset should not exists", p)
		}
	}
}

func (suite *CockroachDatabaseSuite) TestUpdate() {
	as := suite.Require()

	// Given a created asset
	a1 := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a1))

	// When we make some modification to it
	a2 := a1
	// TODO: Should be changes here

	// And we update it
	as.NoError(suite.db.UpdateAssets(context.Background(), a2))

	// Then the pair is updated
	rp, err := suite.db.ReadAssets(context.Background(), a2.Symbol)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(a2, rp[0])
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	as := suite.Require()

	// Given a created asset
	a := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a))

	// When we delete it
	as.NoError(suite.db.DeleteAssets(context.Background(), a))

	// Then we can't read it anymore
	assets, err := suite.db.ReadAssets(context.Background(), a.Symbol)

	// Then there is an error
	as.NoError(err)
	as.Len(assets, 0)
}
