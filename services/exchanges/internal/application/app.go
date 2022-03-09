package application

import (
	"github.com/cryptellation/cryptellation/services/exchanges/internal/application/commands"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ReadAndSyncIfExpiredExchanges commands.ReadAndSyncIfExpiredExchangesHandler
}

type Queries struct {
}
