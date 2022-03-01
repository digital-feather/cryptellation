package application

import (
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePairs commands.CreatePairsHandler
}

type Queries struct {
	ReadPairs queries.ReadPairsHandler
}
