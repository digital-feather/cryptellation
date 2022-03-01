package service

import (
	"github.com/cryptellation/cryptellation/pkg/client"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/assets"
	assetsgrpc "github.com/cryptellation/cryptellation/services/pairs/internal/adapters/assets/grpc"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db/cockroach"
	app "github.com/cryptellation/cryptellation/services/pairs/internal/application"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application/queries"
)

func NewApplication() (app.Application, func(), error) {
	assetsClient, closeAssetsClient, err := client.NewAssetsGrpcClient()
	if err != nil {
		return app.Application{}, func() {}, err
	}
	assetsGrpc := assetsgrpc.NewAssetsGrpc(assetsClient)

	app := newApplication(assetsGrpc)

	closeApplication := func() {
		closeAssetsClient()
	}

	return app, closeApplication, nil
}

func NewMockApplication() app.Application {
	return newApplication(AssetsServiceMock{})
}

func newApplication(assetsService assets.Port) app.Application {
	repository, err := cockroach.New()
	if err != nil {
		return app.Application{}
	}

	return app.Application{
		Commands: app.Commands{
			CreatePairs: commands.NewCreatePairsHandler(repository, assetsService),
		},
		Queries: app.Queries{
			ReadPairs: queries.NewReadPairsHandler(repository),
		},
	}
}
