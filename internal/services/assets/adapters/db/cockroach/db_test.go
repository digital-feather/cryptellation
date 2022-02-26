package cockroach

import (
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/services/assets/adapters/db/cockroach/entities"
	"github.com/cryptellation/cryptellation/internal/utils/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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

	var config cockroachdb.Config
	suite.Require().NoError(config.Load().Validate())

	db, err := New(config)
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *CockroachDatabaseSuite) AfterTest(suiteName, testName string) {
	suite.db.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entities.Asset{})
}

func (suite *CockroachDatabaseSuite) TestNewWithURIError() {
	var err error
	_, err = New(cockroachdb.Config{})
	suite.Error(err)
}

func (suite *CockroachDatabaseSuite) TestCreateRead() {
	as := suite.Require()

	a := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.Create(a))
	rp, err := suite.db.Read(a.Symbol)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(a, rp[0])
}

func (suite *CockroachDatabaseSuite) TestCreateReadInexistant() {
	as := suite.Require()

	a := asset.Asset{Symbol: "ETH"}
	as.NoError(suite.db.Create(a))

	a2 := asset.Asset{Symbol: "BTC"}
	_, err := suite.db.Read(a2.Symbol)
	as.Error(err)
}

func (suite *CockroachDatabaseSuite) TestReadAll() {
	as := suite.Require()

	a1 := asset.Asset{Symbol: "ETH"}
	suite.NoError(suite.db.Create(a1))
	a2 := asset.Asset{Symbol: "FTM"}
	suite.NoError(suite.db.Create(a2))
	a3 := asset.Asset{Symbol: "DAI"}
	suite.NoError(suite.db.Create(a3))

	ps, err := suite.db.Read()
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
	as.NoError(suite.db.Create(a1))
	a2 := a1
	// TODO: Should be changes here
	as.NoError(suite.db.Update(a2))
	rp, err := suite.db.Read(a2.Symbol)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(a2, rp[0])
}

func (suite *CockroachDatabaseSuite) TestDelete() {
	a := asset.Asset{Symbol: "ETH"}
	suite.NoError(suite.db.Create(a))
	suite.NoError(suite.db.Delete(a))
	_, err := suite.db.Read(a.Symbol)
	suite.Error(err)
}
