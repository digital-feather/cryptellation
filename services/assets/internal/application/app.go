package application

import (
	"github.com/cryptellation/cryptellation/services/assets/internal/application/commands"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateAssets commands.CreateAssetsHandler
}

type Queries struct {
	ReadAssets queries.ReadAssetsHandler
}
