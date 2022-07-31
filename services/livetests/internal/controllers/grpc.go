package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/internal/go/controllers/grpc/genproto/livetests"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/event"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
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

func (g GrpcController) SubscribeToLivetestEvents(ctx context.Context, req *livetests.SubscribeToLivetestEventsRequest) (*livetests.SubscribeToLivetestEventsResponse, error) {
	err := g.application.Commands.Livetest.SubscribeToEvents.Handle(ctx, uint(req.Id), req.ExchangeName, req.PairSymbol)
	return &livetests.SubscribeToLivetestEventsResponse{}, err
}

func (g GrpcController) ListenLivetest(srv livetests.LivetestsService_ListenLivetestServer) error {
	// var eventChan <-chan event.Event
	// var err error

	// ctx := srv.Context()
	// firstRequest := true
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		// Exit if context is done or continue
	// 		return ctx.Err()
	// 	case livetestId, ok := <-reqChan:
	// 		if !ok {
	// 			return nil
	// 		}

	// 		// If if it's the first request, then listen to events
	// 		if firstRequest {
	// 			eventChan, err = g.application.Queries.Livetest.ListenEvents.Handle(ctx, livetestId)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			firstRequest = false
	// 		}

	// 		if err = g.application.Commands.Livetest.Advance.Handle(ctx, uint(livetestId)); err != nil {
	// 			return err
	// 		}
	// 	case event, ok := <-eventChan:
	// 		if !ok {
	// 			return nil
	// 		}

	// 		if err := sendEventResponse(srv, event); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func sendEventResponse(srv livetests.LivetestsService_ListenLivetestServer, evt event.Event) error {
	content, err := json.Marshal(evt.Content)
	if err != nil {
		return fmt.Errorf("marshaling event content: %w", err)
	}

	return srv.Send(&livetests.LivetestEventResponse{
		Type:    evt.Type.String(),
		Time:    evt.Time.Format(time.RFC3339),
		Content: string(content),
	})
}
