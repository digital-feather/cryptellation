package binance

import "strings"

func BinanceSymbol(pairSymbol string) string {
	return strings.ReplaceAll(pairSymbol, "-", "")
}
