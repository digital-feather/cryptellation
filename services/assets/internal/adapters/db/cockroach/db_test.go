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
	os.Setenv("COCKROACHDB_DATABASE", "assets")

	db, err, _ := New()
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
	_, err, _ = New()
	suite.Error(err)
}

func (suite *CockroachDatabaseSuite) TestCreateRead() {
	as := suite.Require()

	a := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a))
	rp, err := suite.db.ReadAssets(context.Background(), a.Symbol)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(a, rp[0])
}

func (suite *CockroachDatabaseSuite) TestCreateReadInexistant() {
	as := suite.Require()

	a := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a))

	a2 := asset.Asset{Symbol: "BTC"}
	_, err := suite.db.ReadAssets(context.Background(), a2.Symbol)
	as.Error(err)
}

func (suite *CockroachDatabaseSuite) TestReadAll() {
	as := suite.Require()

	a1 := asset.Asset{Symbol: "ETH"}
	suite.NoError(suite.db.CreateAssets(context.Background(), a1))
	a2 := asset.Asset{Symbol: "FTM"}
	suite.NoError(suite.db.CreateAssets(context.Background(), a2))
	a3 := asset.Asset{Symbol: "DAI"}
	suite.NoError(suite.db.CreateAssets(context.Background(), a3))

	ps, err := suite.db.ReadAssets(context.Background())
	as.NoError(err)
	as.Len(ps, 3)

	for _, p := range ps {
		if p.Symbol != a1.Symbol && p.Symbol != a2.Symbol && p.Symbol != a3.Symbol {
			as.Fail("This asset should not exists", p)
		}
	}
}

func (suite *CockroachDatabaseSuite) TestUpdate() {
	as := suite.Require()

	a1 := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.CreateAssets(context.Background(), a1))
	a2 := a1
	// TODO: Should be changes here
	as.NoError(suite.db.UpdateAssets(context.Background(), a2))
	rp, err := suite.db.ReadAssets(context.Background(), a2.Symbol)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(a2, rp[0])
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	a := asset.Asset{Symbol: "ETH"}
	suite.NoError(suite.db.CreateAssets(context.Background(), a))
	suite.NoError(suite.db.DeleteAssets(context.Background(), a))
	_, err := suite.db.ReadAssets(context.Background(), a.Symbol)
	suite.Error(err)
}
