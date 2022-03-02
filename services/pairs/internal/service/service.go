package service

import (
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db/cockroach"
	app "github.com/cryptellation/cryptellation/services/pairs/internal/application"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/queries"
)

func NewApplication() (app.Application, error) {
	repository, err := cockroach.New()
	if err != nil {
		return app.Application{}, err
	}

	return app.Application{
		Commands: app.Commands{
			CreatePairs: commands.NewCreatePairsHandler(repository),
		},
		Queries: app.Queries{
			ReadPairs: queries.NewReadPairsHandler(repository),
		},
	}, nil
}
