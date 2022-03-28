package service

import (
	client "github.com/cryptellation/cryptellation/clients/go"
	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	pubsubRedis "github.com/cryptellation/cryptellation/services/backtests/internal/adapters/pubsub/redis"
	vdbRedis "github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb/redis"
	app "github.com/cryptellation/cryptellation/services/backtests/internal/application"
	cmdBacktest "github.com/cryptellation/cryptellation/services/backtests/internal/application/commands/backtest"
	queriesBacktest "github.com/cryptellation/cryptellation/services/backtests/internal/application/queries/backtest"
)

func NewApplication() (app.Application, func(), error) {
	client, closeClient, err := client.NewCandlesticksGrpcClient()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	app, closeApp, err := newApplication(client)

	return app, func() {
		closeApp()
		closeClient()
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication(MockedCandlesticksClient{})
}

func newApplication(client candlesticks.CandlesticksServiceClient) (app.Application, func(), error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	ps, err := pubsubRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			Backtest: app.BacktestCommands{
				Advance:           cmdBacktest.NewAdvanceHandler(repository, ps, client),
				Create:            cmdBacktest.NewCreateHandler(repository),
				SubscribeToEvents: cmdBacktest.NewSubscribeToEventsHandler(repository),
			},
		},
		Queries: app.Queries{
			Backtest: app.BacktestQueries{
				ListenEvents: queriesBacktest.NewListenEventsHandler(ps),
			},
		},
	}, func() {}, nil
}
