package controllers

import (
	"context"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/ticks"
	app "github.com/digital-feather/cryptellation/services/ticks/internal/application"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) ListenSymbol(req *ticks.ListenSymbolRequest, srv ticks.TicksService_ListenSymbolServer) error {
	ctx := srv.Context()

	// Start listening before registration to avoid missing ticks
	ticksChanRecv, err := g.application.Queries.ListenSymbol.Handle(ctx, req.Exchange, req.PairSymbol)
	if err != nil {
		return err
	}

	err = g.application.Commands.RegisterSymbolListener.Handle(ctx, req.Exchange, req.PairSymbol)
	if err != nil {
		return err
	}

	loopErr := loopOverNewTicks(ctx, srv, ticksChanRecv)
	unregisterErr := g.application.Commands.UnregisterSymbolListener.Handle(context.Background(), req.Exchange, req.PairSymbol)

	if loopErr == nil {
		return unregisterErr
	}

	log.Println(unregisterErr)
	return loopErr
}

func loopOverNewTicks(ctx context.Context, srv ticks.TicksService_ListenSymbolServer, ticksChanRecv <-chan tick.Tick) error {
	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		t, ok := <-ticksChanRecv
		if !ok {
			return nil
		}

		if err := srv.Send(toGrpcTick(t)); err != nil {
			return err
		}
	}
}

func toGrpcTick(t tick.Tick) *ticks.Tick {
	return &ticks.Tick{
		Time:       t.Time.Format(time.RFC3339Nano),
		Exchange:   t.Exchange,
		PairSymbol: t.PairSymbol,
		Price:      float32(t.Price),
	}
}
