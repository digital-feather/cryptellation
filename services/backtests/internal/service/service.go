package service

import (
	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/genproto/candlesticks"
	pubsubRedis "github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub/redis"
	vdbRedis "github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/backtests/internal/application"
	cmdBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/commands/backtest"
	queriesBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/queries/backtest"
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
				CreateOrder:       cmdBacktest.NewCreateOrderHandler(repository, client),
				SubscribeToEvents: cmdBacktest.NewSubscribeToEventsHandler(repository),
			},
		},
		Queries: app.Queries{
			Backtest: app.BacktestQueries{
				GetAccounts:  queriesBacktest.NewGetAccounts(repository),
				GetOrders:    queriesBacktest.NewGetOrders(repository),
				ListenEvents: queriesBacktest.NewListenEventsHandler(ps),
			},
		},
	}, func() {}, nil
}
