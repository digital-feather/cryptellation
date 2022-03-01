package domain

import (
	"testing"

	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/stretchr/testify/suite"
)

func TestCreatePairsSuite(t *testing.T) {
	suite.Run(t, new(CreatePairsSuite))
}

type CreatePairsSuite struct {
	suite.Suite
}

func (suite *CreatePairsSuite) TestCreatePairsOk() {
	a1 := asset.Asset{
		Symbol: "ETH",
	}
	a2 := asset.Asset{
		Symbol: "USDC",
	}
	p1 := pair.New(a1.Symbol, a2.Symbol)

	a3 := asset.Asset{
		Symbol: "BTC",
	}
	a4 := asset.Asset{
		Symbol: "USDT",
	}
	p2 := pair.New(a3.Symbol, a4.Symbol)

	pairs, err := CreatePairs([]pair.Pair{p1, p2}, []asset.Asset{a1, a2, a3, a4})
	suite.Require().NoError(err)
	suite.Equal(p1, pairs[0])
	suite.Equal(p2, pairs[1])
}

func (suite *CreatePairsSuite) TestCreatePairsInexistantAsset() {
	a1 := asset.Asset{
		Symbol: "ETH",
	}
	a2 := asset.Asset{
		Symbol: "USDC",
	}
	p1 := pair.New(a1.Symbol, a2.Symbol)

	a3 := asset.Asset{
		Symbol: "BTC",
	}
	a4 := asset.Asset{
		Symbol: "USDT",
	}
	p2 := pair.New(a3.Symbol, a4.Symbol)

	_, err := CreatePairs([]pair.Pair{p1, p2}, []asset.Asset{a1, a2, a3})
	suite.Require().Error(err)
}
