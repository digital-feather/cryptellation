package controllers

import (
	"context"
	"log"

	"github.com/cryptellation/cryptellation/internal/genproto/pairs"
	"github.com/cryptellation/cryptellation/pkg/types/pair"
	app "github.com/cryptellation/cryptellation/services/pairs/internal/application"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) CreatePairs(ctx context.Context, req *pairs.CreatePairsRequest) (*pairs.CreatePairsResponse, error) {
	if err := g.application.Commands.CreatePairs.Handle(ctx, fromGrpcPairs(req.Pairs)); err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pairs.CreatePairsResponse{}, nil
}

func (g GrpcController) ReadPairs(ctx context.Context, req *pairs.ReadPairsRequest) (*pairs.ReadPairsResponse, error) {
	list, err := g.application.Queries.ReadPairs.Handle(ctx, req.Symbols)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pairs.ReadPairsResponse{
		Pairs: toGrpcPairs(list),
	}, nil
}

func fromGrpcPairs(gpairs []*pairs.Pair) []pair.Pair {
	ps := make([]pair.Pair, len(gpairs))
	for i, p := range gpairs {
		ps[i] = pair.Pair{
			BaseAssetSymbol:  p.BaseAssetSymbol,
			QuoteAssetSymbol: p.QuoteAssetSymbol,
		}
	}
	return ps
}

func toGrpcPairs(ps []pair.Pair) []*pairs.Pair {
	gpairs := make([]*pairs.Pair, len(ps))
	for i, p := range ps {
		gpairs[i] = &pairs.Pair{
			Symbol:           p.Symbol(),
			BaseAssetSymbol:  p.BaseAssetSymbol,
			QuoteAssetSymbol: p.QuoteAssetSymbol,
		}
	}
	return gpairs
}
