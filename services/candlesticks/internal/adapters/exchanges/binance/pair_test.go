package binance

import (
	"testing"
)

func TestBinanceSymbol(t *testing.T) {
	if BinanceSymbol("ETH-USDC") != "ETHUSDC" {
		t.Error("Not equal")
	}
}
