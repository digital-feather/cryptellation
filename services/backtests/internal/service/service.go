package service

import (
	"log"

	pubsubRedis "github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub/redis"
	vdbRedis "github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/backtests/internal/application"
	cmdBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/commands/backtest"
	queriesBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/queries/backtest"
	candlesticksGrpc "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	candlesticksProto "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
)

func NewApplication() (app.Application, func(), error) {
	csClient, closeCsClient, err := candlesticksGrpc.Newclient()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	app, closeApp, err := newApplication(csClient)

	return app, func() {
		closeApp()
		if err := closeCsClient(); err != nil {
			log.Println("error when closing candlestick client:", err)
		}
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication(MockedCandlesticksClient{})
}

func newApplication(client candlesticksProto.CandlesticksServiceClient) (app.Application, func(), error) {
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
