package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cryptellation/cryptellation/internal/genproto/backtests"
	app "github.com/cryptellation/cryptellation/services/backtests/internal/application"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/account"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/event"
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

	id, err := g.application.Commands.Backtest.Create.Handle(ctx, newPayload)
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
		Accounts:              acc,
		StartTime:             st,
		EndTime:               et,
		DurationBetweenEvents: tbe,
	}, nil
}

func (g GrpcController) ListenBacktest(req *backtests.ListenBacktestRequest, srv backtests.BacktestsService_ListenBacktestServer) error {
	ctx := srv.Context()

	eventsChanRecv, err := g.application.Queries.Backtest.ListenEvents.Handle(ctx, req.Id)
	if err != nil {
		return err
	}

	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		event, ok := <-eventsChanRecv
		if !ok {
			return nil
		}

		if err := sendEvent(srv, event); err != nil {
			log.Println("error when sending event:", err)
			return nil
		}
	}
}

func sendEvent(srv backtests.BacktestsService_ListenBacktestServer, evt event.Interface) error {
	content, err := json.Marshal(evt.GetContent())
	if err != nil {
		return xerrors.Errorf("marshaling event content: %w", err)
	}

	return srv.Send(&backtests.Event{
		Type:    evt.GetType().String(),
		Time:    evt.GetTime().Format(time.RFC3339),
		Content: string(content),
	})
}

func (g GrpcController) SubscribeToBacktestEvents(ctx context.Context, req *backtests.SubscribeToBacktestEventsRequest) (*backtests.SubscribeToBacktestEventsResponse, error) {
	err := g.application.Commands.Backtest.SubscribeToEvents.Handle(ctx, uint(req.Id), req.Exchange, req.PairSymbol)
	return &backtests.SubscribeToBacktestEventsResponse{}, err
}

func (g GrpcController) AdvanceBacktest(ctx context.Context, req *backtests.AdvanceBacktestRequest) (*backtests.AdvanceBacktestResponse, error) {
	finished, err := g.application.Commands.Backtest.Advance.Handle(ctx, uint(req.Id))
	return &backtests.AdvanceBacktestResponse{
		Finished: finished,
	}, err
}
