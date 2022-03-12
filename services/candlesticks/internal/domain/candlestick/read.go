package candlestick

import (
	"time"

	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/period"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func AreMissing(cl *List, start, end time.Time, limit uint) bool {
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

func GetDownloadStartEndTimes(cl *List, start, end time.Time) (time.Time, time.Time) {
	t, _, exists := cl.Last()
	if exists && !cl.HasUncomplete() {
		start = t.Add(cl.Period().Duration())
	}

	qty := int(cl.Period().CountBetweenTimes(start, end)) + 1
	if qty < MinimalRetrievedMissingCandlesticks {
		d := cl.Period().Duration() * time.Duration(MinimalRetrievedMissingCandlesticks-qty)
		end = end.Add(d)
	}

	return start, end
}

func ProcessRequestedStartEndTimes(per period.Symbol, start, end *time.Time) (time.Time, time.Time) {
	var nstart, nend time.Time

	defaultDuration := per.Duration() * 500
	if end == nil {
		if start == nil {
			nend = time.Now()
		} else {
			nend = start.Add(defaultDuration)
		}
	} else {
		nend = *end
	}

	if start == nil {
		nstart = nend.Add(-defaultDuration)
	} else {
		nstart = *start
	}

	nstart = per.RoundTime(nstart)
	nend = per.RoundTime(nend)

	return nstart, nend
}
