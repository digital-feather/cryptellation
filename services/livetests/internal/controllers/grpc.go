package controllers

import (
	"context"

	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/livetests"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/account"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) CreateLivetest(ctx context.Context, req *livetests.CreateLivetestRequest) (*livetests.CreateLivetestResponse, error) {
	newPayload, err := fromCreateLivetestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.application.Commands.Livetest.Create.Handle(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &livetests.CreateLivetestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateLivetestRequest(req *livetests.CreateLivetestRequest) (livetest.NewPayload, error) {
	acc := make(map[string]account.Account, len(req.Accounts))
	for exch, v := range req.Accounts {
		balances := make(map[string]float64, len(v.Assets))
		for asset, qty := range v.Assets {
			balances[asset] = float64(qty)
		}

		acc[exch] = account.Account{
			Balances: balances,
		}
	}

	return livetest.NewPayload{
		Accounts: acc,
	}, nil
}
