package utils

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

var (
	ErrInvalidPairSymbol = errors.New("invalid pair symbol")
)

func FormatPairSymbol(baseSymbol, quoteSymbol string) string {
	return fmt.Sprintf("%s-%s", baseSymbol, quoteSymbol)
}

func ParsePairSymbol(symbol string) (baseSymbol, quoteSymbol string, err error) {
	split := strings.Split(symbol, "-")
	if len(split) != 2 {
		return "", "", xerrors.Errorf("error parsing pair symbol %q: %w", symbol, ErrInvalidPairSymbol)
	}

	return split[0], split[1], nil
}
