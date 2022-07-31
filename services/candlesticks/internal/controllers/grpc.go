package controllers

import (
	"context"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/internal/go/controllers/grpc/genproto/candlesticks"
	app "github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/commands"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/period"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) ReadCandlesticks(ctx context.Context, req *candlesticks.ReadCandlesticksRequest) (*candlesticks.ReadCandlesticksResponse, error) {
	payload, err := fromReadCandlesticksRequest(req)
	if err != nil {
		return nil, err
	}

	list, err := g.application.Commands.CachedReadCandlesticks.Handle(ctx, payload)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &candlesticks.ReadCandlesticksResponse{
		Candlesticks: toGrpcCandlesticks(list),
	}, nil
}

func fromReadCandlesticksRequest(req *candlesticks.ReadCandlesticksRequest) (commands.CachedReadCandlesticksPayload, error) {
	per, err := period.FromString(req.PeriodSymbol)
	if err != nil {
		return commands.CachedReadCandlesticksPayload{}, err
	}

	payload := commands.CachedReadCandlesticksPayload{
		ExchangeName: req.ExchangeName,
		PairSymbol:   req.PairSymbol,
		Period:       per,
		Limit:        uint(req.Limit),
	}

	if req.Start != "" {
		start, err := time.Parse(time.RFC3339Nano, req.Start)
		if err != nil {
			return commands.CachedReadCandlesticksPayload{}, err
		}
		payload.Start = &start
	}

	if req.End != "" {
		end, err := time.Parse(time.RFC3339Nano, req.End)
		if err != nil {
			return commands.CachedReadCandlesticksPayload{}, err
		}
		payload.End = &end
	}

	return payload, nil
}

func toGrpcCandlesticks(cl *candlestick.List) []*candlesticks.Candlestick {
	gcandlesticks := make([]*candlesticks.Candlestick, 0, cl.Len())
	_ = cl.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		gcandlesticks = append(gcandlesticks, &candlesticks.Candlestick{
			Time:   t.Format(time.RFC3339Nano),
			Open:   float32(cs.Open),
			High:   float32(cs.High),
			Low:    float32(cs.Low),
			Close:  float32(cs.Close),
			Volume: float32(cs.Volume),
		})

		return false, nil
	})
	return gcandlesticks
}
