package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/internal/genproto/exchanges"
	app "github.com/digital-feather/cryptellation/services/exchanges/internal/application"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) ReadExchanges(ctx context.Context, req *exchanges.ReadExchangesRequest) (*exchanges.ReadExchangesResponse, error) {
	list, err := g.application.Commands.CachedReadExchanges.Handle(ctx, nil, req.Names...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &exchanges.ReadExchangesResponse{
		Exchanges: toGrpcExchanges(list),
	}, nil
}

func fromGrpcExchanges(gexchanges []*exchanges.Exchange) ([]exchange.Exchange, error) {
	ps := make([]exchange.Exchange, len(gexchanges))
	for i, p := range gexchanges {
		lastSyncTime, err := time.Parse(time.RFC3339, p.LastSyncTime)
		if err != nil {
			return nil, xerrors.New(fmt.Sprintf("invalid time %q for %q exchange", p.LastSyncTime, p.Name))
		}

		ps[i] = exchange.Exchange{
			Name:           p.Name,
			PairsSymbols:   p.Pairs,
			PeriodsSymbols: p.Periods,
			Fees:           float64(p.Fees),
			LastSyncTime:   lastSyncTime,
		}
	}
	return ps, nil
}

func toGrpcExchanges(ps []exchange.Exchange) []*exchanges.Exchange {
	gexchanges := make([]*exchanges.Exchange, len(ps))
	for i, p := range ps {
		gexchanges[i] = &exchanges.Exchange{
			Name:         p.Name,
			Pairs:        p.PairsSymbols,
			Periods:      p.PeriodsSymbols,
			Fees:         float32(p.Fees),
			LastSyncTime: p.LastSyncTime.Format(time.RFC3339),
		}
	}
	return gexchanges
}
