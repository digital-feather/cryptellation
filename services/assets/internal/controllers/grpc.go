package controllers

import (
	"context"
	"log"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	app "github.com/cryptellation/cryptellation/services/assets/internal/application"
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
	if err := g.application.Commands.CreateAssets.Handle(ctx, fromGrpcAssets(req.Assets)); err != nil {
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

	return &assets.ReadAssetsResponse{
		Assets: toGrpcAssets(list),
	}, nil
}

func fromGrpcAssets(gassets []*assets.Asset) []asset.Asset {
	ps := make([]asset.Asset, len(gassets))
	for i, p := range gassets {
		ps[i] = asset.Asset{
			Symbol: p.Symbol,
		}
	}
	return ps
}

func toGrpcAssets(ps []asset.Asset) []*assets.Asset {
	gassets := make([]*assets.Asset, len(ps))
	for i, p := range ps {
		gassets[i] = &assets.Asset{
			Symbol: p.Symbol,
		}
	}
	return gassets
}
