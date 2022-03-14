package controllers

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/internal/genproto/backtests"
	app "github.com/cryptellation/cryptellation/services/backtests/internal/application"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/account"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/backtest"
	"golang.org/x/xerrors"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) CreateBacktest(ctx context.Context, req *backtests.CreateBacktestRequest) (*backtests.CreateBacktestResponse, error) {
	newPayload, err := fromCreateBacktestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.application.Commands.CreateBacktest.Handle(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &backtests.CreateBacktestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateBacktestRequest(req *backtests.CreateBacktestRequest) (backtest.NewPayload, error) {
	st, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return backtest.NewPayload{}, xerrors.Errorf("error when parsing start_time: %w", err)
	}

	var et *time.Time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return backtest.NewPayload{}, xerrors.Errorf("error when parsing start_time: %w", err)
		}
		et = &t
	}

	var tbe *time.Duration
	if req.SecondsBetweenPriceEvents > 0 {
		d := time.Duration(req.SecondsBetweenPriceEvents) * time.Second
		tbe = &d
	}

	acc := make(map[string]account.Account, len(req.Accounts))
	for _, v := range req.Accounts {
		balances := make(map[string]float64, len(v.Assets))
		for _, a := range v.Assets {
			balances[a.AssetName] = float64(a.Quantity)
		}

		acc[v.ExchangeName] = account.Account{
			Balances: balances,
		}
	}

	return backtest.NewPayload{
		Accounts:          acc,
		StartTime:         st,
		EndTime:           et,
		TimeBetweenEvents: tbe,
	}, nil
}
