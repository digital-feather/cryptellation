package service

import (
	pubsubRedis "github.com/cryptellation/cryptellation/services/backtests/internal/adapters/pubsub/redis"
	vdbRedis "github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb/redis"
	app "github.com/cryptellation/cryptellation/services/backtests/internal/application"
	"github.com/cryptellation/cryptellation/services/backtests/internal/application/commands"
)

func NewApplication() (app.Application, error) {
	repository, err := vdbRedis.New()
	if err != nil {
		return app.Application{}, err
	}

	ps, err := pubsubRedis.New()
	if err != nil {
		return app.Application{}, err
	}

	return app.Application{
		Commands: app.Commands{
			CreateBacktest: commands.NewCreateBacktestHandler(repository, ps),
		},
	}, nil
}
