package exchanges

import "github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"

type Port interface {
	ListenSymbol(symbol string) (chan tick.Tick, chan struct{}, error)
}
