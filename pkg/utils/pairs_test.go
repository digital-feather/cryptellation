package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

type UtilsSuite struct {
	suite.Suite
}

func (suite *UtilsSuite) TestFormatPair() {
	cases := []struct {
		BaseSymbol  string
		QuoteSymbol string
		PairSymbol  string
	}{
		{
			BaseSymbol:  "ETH",
			QuoteSymbol: "BTC",
			PairSymbol:  "ETH-BTC",
		},
	}

	for i, c := range cases {
		symbol := FormatPairSymbol(c.BaseSymbol, c.QuoteSymbol)
		suite.Require().Equal(c.PairSymbol, symbol, i)
	}
}

func (suite *UtilsSuite) TestParsePair() {
	cases := []struct {
		PairSymbol  string
		BaseSymbol  string
		QuoteSymbol string
		Error       bool
	}{
		{
			PairSymbol:  "ETH-BTC",
			BaseSymbol:  "ETH",
			QuoteSymbol: "BTC",
			Error:       false,
		}, {
			PairSymbol:  "",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       true,
		}, {
			PairSymbol:  "-",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       false,
		}, {
			PairSymbol:  "--",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       true,
		},
	}

	for i, c := range cases {
		baseSymbol, quoteSymbol, err := ParsePairSymbol(c.PairSymbol)
		suite.Require().Equal(c.BaseSymbol, baseSymbol, i)
		suite.Require().Equal(c.QuoteSymbol, quoteSymbol, i)
		if c.Error {
			suite.Require().Error(err, i)
		} else {
			suite.Require().NoError(err, i)
		}
	}
}
