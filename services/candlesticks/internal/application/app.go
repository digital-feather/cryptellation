package application

import "github.com/digital-feather/cryptellation/services/candlesticks/internal/application/commands"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CachedReadCandlesticks commands.CachedReadCandlesticksHandler
}

type Queries struct {
}
