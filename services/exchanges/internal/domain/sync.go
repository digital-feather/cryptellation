package domain

import (
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/exchange"
)

const DefaultExpirationDuration = time.Hour

func GetExpiredExchangesNames(
	expectedExchanges []string,
	exchangesFromDB []exchange.Exchange,
	expirationDuration *time.Duration,
) (toSync []string, err error) {
	mappedExchanges := exchange.ArrayToMap(exchangesFromDB)

	if expirationDuration == nil {
		d := DefaultExpirationDuration
		expirationDuration = &d
	}

	toSync = make([]string, 0, len(expectedExchanges))
	for _, name := range expectedExchanges {
		exch, ok := mappedExchanges[name]
		if ok && time.Now().Before(exch.LastSyncTime.Add(*expirationDuration)) {
			continue
		}

		toSync = append(toSync, name)
	}

	return toSync, err
}
