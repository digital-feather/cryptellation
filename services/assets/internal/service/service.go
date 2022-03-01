package service

import (
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db/cockroach"
	app "github.com/cryptellation/cryptellation/services/assets/internal/application"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/queries"
)

func NewApplication() (app.Application, error) {
	repository, err := cockroach.New()
	if err != nil {
		return app.Application{}, err
	}

	a := app.Application{
		Commands: app.Commands{
			CreateAssets: commands.NewCreateAssetHandler(repository),
		},
		Queries: app.Queries{
			ReadAssets: queries.NewReadAssetsHandler(repository),
		},
	}

	return a, nil
}
