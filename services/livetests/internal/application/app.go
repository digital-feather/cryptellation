package application

import (
	cmdLivetest "github.com/digital-feather/cryptellation/services/livetests/internal/application/commands/livetest"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type LivetestCommands struct {
	Create cmdLivetest.CreateHandler
}

type Commands struct {
	Livetest LivetestCommands
}

type LivetestQueries struct {
}

type Queries struct {
	Livetest LivetestQueries
}
