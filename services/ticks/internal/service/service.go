package service

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges/binance"
	pubsubRedis "github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub/redis"
	vdbRedis "github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/ticks/internal/application"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/commands"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/queries"
)

func NewApplication() (app.Application, func(), error) {
	exchanges, err := instanciateExchanges()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	app, closeApp, err := newApplication(exchanges)

	return app, func() {
		closeApp()
	}, err
}

func instanciateExchanges() (map[string]exchanges.Port, error) {
	var err error
	exchanges := make(map[string]exchanges.Port)

	exchanges[binance.Name], err = binance.New()
	if err != nil {
		return nil, err
	}

	return exchanges, nil
}

func NewMockedApplication() (app.Application, func(), error) {
	services := map[string]exchanges.Port{
		MockExchangeName: &MockExchangeService{},
	}

	return newApplication(services)
}

func newApplication(exchanges map[string]exchanges.Port) (app.Application, func(), error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	if err := repository.ClearSymbolListenersCount(context.TODO()); err != nil {
		return app.Application{}, func() {}, err
	}

	ps, err := pubsubRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			RegisterSymbolListener:   commands.NewRegisterSymbolListener(ps, repository, exchanges),
			UnregisterSymbolListener: commands.NewUnregisterSymbolListener(repository),
		},
		Queries: app.Queries{
			ListenSymbol: queries.NewListenSymbolsHandler(ps),
		},
	}, func() {}, nil
}
