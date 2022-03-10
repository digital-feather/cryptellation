package domain

import (
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func AreCsMissing(cl *candlestick.List, start, end time.Time, limit uint) bool {
	expectedCount := int(cl.Period().CountBetweenTimes(start, end)) + 1
	qty := cl.Len()

	if qty < expectedCount && (limit == 0 || uint(qty) < limit) {
		return true
	}

	if cl.HasUncomplete() {
		return true
	}

	return false
}

func GetRequestedCandlesticksFromList(cl *candlestick.List, start, end time.Time, limit uint) *candlestick.List {
	ecl := cl.Extract(start, end)
	if limit == 0 || ecl.Len() < int(limit) {
		return ecl
	}

	return ecl.FirstN(limit)
}

func MinimalCandlesticksEndTimeForDownload(cl *candlestick.List, start, end time.Time) time.Time {
	qty := int(cl.Period().CountBetweenTimes(start, end)) + 1
	if qty < MinimalRetrievedMissingCandlesticks {
		d := cl.Period().Duration() * time.Duration(MinimalRetrievedMissingCandlesticks-qty)
		end = end.Add(d)
	}

	return end
}
