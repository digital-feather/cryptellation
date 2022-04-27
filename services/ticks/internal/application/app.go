package application

import (
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/commands"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterSymbolListener   commands.RegisterSymbolListenerHandler
	UnregisterSymbolListener commands.UnregisterSymbolListenerHandler
}

type Queries struct {
	ListenSymbol queries.ListenSymbolsHandler
}
