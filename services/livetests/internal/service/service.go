package service

import (
	vdbRedis "github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	cmdLivetest "github.com/digital-feather/cryptellation/services/livetests/internal/application/commands/livetest"
)

func NewApplication() (app.Application, func(), error) {
	app, closeApp, err := newApplication()
	return app, func() {
		closeApp()
	}, err
}

func NewMockedApplication() (app.Application, func(), error) {
	return newApplication()
}

func newApplication() (app.Application, func(), error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, func() {}, err
	}

	return app.Application{
		Commands: app.Commands{
			Livetest: app.LivetestCommands{
				Create: cmdLivetest.NewCreateHandler(repository),
			},
		},
		Queries: app.Queries{
			Livetest: app.LivetestQueries{},
		},
	}, func() {}, nil
}
