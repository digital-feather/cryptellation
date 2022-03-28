package application

import (
	cmdBacktest "github.com/cryptellation/cryptellation/services/backtests/internal/application/commands/backtest"
	queriesBacktest "github.com/cryptellation/cryptellation/services/backtests/internal/application/queries/backtest"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type BacktestCommands struct {
	Advance           cmdBacktest.AdvanceHandler
	Create            cmdBacktest.CreateHandler
	SubscribeToEvents cmdBacktest.SubscribeToEventsHandler
}

type Commands struct {
	Backtest BacktestCommands
}

type BacktestQueries struct {
	ListenEvents queriesBacktest.ListenEventsHandler
}

type Queries struct {
	Backtest BacktestQueries
}
