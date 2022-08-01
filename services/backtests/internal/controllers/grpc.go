package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	app "github.com/digital-feather/cryptellation/services/backtests/internal/application"
	cmdBacktest "github.com/digital-feather/cryptellation/services/backtests/internal/application/commands/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/order"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) CreateBacktest(ctx context.Context, req *proto.CreateBacktestRequest) (*proto.CreateBacktestResponse, error) {
	newPayload, err := fromCreateBacktestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.application.Commands.Backtest.Create.Handle(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &proto.CreateBacktestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateBacktestRequest(req *proto.CreateBacktestRequest) (backtest.NewPayload, error) {
	st, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return backtest.NewPayload{}, fmt.Errorf("error when parsing start_time: %w", err)
	}

	var et *time.Time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return backtest.NewPayload{}, fmt.Errorf("error when parsing start_time: %w", err)
		}
		et = &t
	}

	var tbe *time.Duration
	if req.SecondsBetweenPriceEvents > 0 {
		d := time.Duration(req.SecondsBetweenPriceEvents) * time.Second
		tbe = &d
	}

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

	return backtest.NewPayload{
		Accounts:              acc,
		StartTime:             st,
		EndTime:               et,
		DurationBetweenEvents: tbe,
	}, nil
}

func (g GrpcController) ListenBacktest(srv proto.BacktestsService_ListenBacktestServer) error {
	var eventChan <-chan event.Event
	var err error

	// Get event requests channel
	reqChan := listenEventRequests(srv)

	ctx := srv.Context()
	firstRequest := true
	for {
		select {
		case <-ctx.Done():
			// Exit if context is done or continue
			return ctx.Err()
		case backtestId, ok := <-reqChan:
			if !ok {
				return nil
			}

			// If if it's the first request, then listen to events
			if firstRequest {
				eventChan, err = g.application.Queries.Backtest.ListenEvents.Handle(ctx, backtestId)
				if err != nil {
					return err
				}
				firstRequest = false
			}

			if err = g.application.Commands.Backtest.Advance.Handle(ctx, uint(backtestId)); err != nil {
				return err
			}
		case event, ok := <-eventChan:
			if !ok {
				return nil
			}

			if err := sendEventResponse(srv, event); err != nil {
				return err
			}
		}
	}
}

func listenEventRequests(srv proto.BacktestsService_ListenBacktestServer) <-chan uint64 {
	requests := make(chan uint64)

	go func() {
		for {
			req, err := srv.Recv()
			if err != nil {
				// return will close stream from server side
				close(requests)
				break
			}

			requests <- req.Id
		}
	}()

	return requests
}

func sendEventResponse(srv proto.BacktestsService_ListenBacktestServer, evt event.Event) error {
	content, err := json.Marshal(evt.Content)
	if err != nil {
		return fmt.Errorf("marshaling event content: %w", err)
	}

	return srv.Send(&proto.BacktestEventResponse{
		Type:    evt.Type.String(),
		Time:    evt.Time.Format(time.RFC3339),
		Content: string(content),
	})
}

func (g GrpcController) SubscribeToBacktestEvents(ctx context.Context, req *proto.SubscribeToBacktestEventsRequest) (*proto.SubscribeToBacktestEventsResponse, error) {
	err := g.application.Commands.Backtest.SubscribeToEvents.Handle(ctx, uint(req.Id), req.ExchangeName, req.PairSymbol)
	return &proto.SubscribeToBacktestEventsResponse{}, err
}

func (g GrpcController) CreateBacktestOrder(ctx context.Context, req *proto.CreateBacktestOrderRequest) (*proto.CreateBacktestOrderResponse, error) {
	payload := cmdBacktest.CreateOrderPayload{
		BacktestId:   uint(req.BacktestId),
		Type:         order.Type(req.Type),
		ExchangeName: req.ExchangeName,
		PairSymbol:   req.PairSymbol,
		Side:         order.Side(req.Side),
		Quantity:     float64(req.Quantity),
	}

	err := g.application.Commands.Backtest.CreateOrder.Handle(ctx, payload)
	return &proto.CreateBacktestOrderResponse{}, err
}

func (g GrpcController) Accounts(ctx context.Context, req *proto.AccountsRequest) (*proto.AccountsResponse, error) {
	accounts, err := g.application.Queries.Backtest.GetAccounts.Handle(ctx, uint(req.BacktestId))
	if err != nil {
		return nil, err
	}

	resp := proto.AccountsResponse{
		Accounts: make(map[string]*proto.Account, len(accounts)),
	}

	for exch, acc := range accounts {
		resp.Accounts[exch] = toGrpcAccount(exch, acc)
	}

	return &resp, nil
}

func toGrpcAccount(exchange string, account account.Account) *proto.Account {
	assets := make(map[string]float32, len(account.Balances))
	for asset, qty := range account.Balances {
		assets[asset] = float32(qty)
	}

	return &proto.Account{
		Assets: assets,
	}
}

func (g GrpcController) Orders(ctx context.Context, req *proto.OrdersRequest) (*proto.OrdersResponse, error) {
	orders, err := g.application.Queries.Backtest.GetOrders.Handle(ctx, uint(req.BacktestId))
	if err != nil {
		return nil, err
	}

	return &proto.OrdersResponse{
		Orders: toGrpcOrders(orders),
	}, nil
}

func toGrpcOrders(orders []order.Order) []*proto.Order {
	formattedOrders := make([]*proto.Order, len(orders))
	for i, o := range orders {
		formattedOrders[i] = &proto.Order{
			Time:         o.Time.Format(time.RFC3339),
			Type:         o.Type.String(),
			ExchangeName: o.ExchangeName,
			PairSymbol:   o.PairSymbol,
			Side:         o.Side.String(),
			Quantity:     float32(o.Quantity),
			Price:        float32(o.Price),
		}
	}
	return formattedOrders
}
