package controllers

import (
	"context"
	"log"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	app "github.com/cryptellation/cryptellation/services/assets/internal/application"
	"github.com/cryptellation/cryptellation/services/assets/internal/application/commands"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) CreateAssets(ctx context.Context, req *assets.CreateAssetsRequest) (*assets.CreateAssetsResponse, error) {
	cmd := commands.CreateAssets{
		Assets: make([]asset.Asset, len(req.Assets)),
	}

	for i, a := range req.Assets {
		cmd.Assets[i] = asset.Asset{
			Symbol: a.Symbol,
		}
	}

	if err := g.application.Commands.CreateAssets.Handle(ctx, cmd); err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &assets.CreateAssetsResponse{}, nil
}

func (g GrpcController) ReadAssets(ctx context.Context, req *assets.ReadAssetsRequest) (*assets.ReadAssetsResponse, error) {
	list, err := g.application.Queries.ReadAssets.Handle(ctx, req.Symbols)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := assets.ReadAssetsResponse{
		Assets: make([]*assets.Asset, len(list)),
	}
	for i, a := range list {
		resp.Assets[i] = &assets.Asset{
			Symbol: a.Symbol,
		}
	}

	return &assets.ReadAssetsResponse{}, nil
}
