package application

import "github.com/cryptellation/cryptellation/services/backtests/internal/application/commands"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBacktest commands.CreateBacktestHandler
}

type Queries struct {
}
