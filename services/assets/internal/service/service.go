package service

import (
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db/cockroach"
	app "github.com/cryptellation/cryptellation/services/assets/internal/application"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/queries"
)

func NewApplication() (app.Application, error, func()) {
	repository, err, closeRepository := cockroach.New()
	if err != nil {
		return app.Application{}, err, func() {}
	}

	a := app.Application{
		Commands: app.Commands{
			CreateAssets: commands.NewCreateAssetHandler(repository),
		},
		Queries: app.Queries{
			ReadAssets: queries.NewReadAssetsHandler(repository),
		},
	}

	closeApplication := func() {
		closeRepository()
	}

	return a, nil, closeApplication
}
