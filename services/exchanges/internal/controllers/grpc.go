package controllers

import (
	"context"
	"log"
	"time"

	app "github.com/digital-feather/cryptellation/services/exchanges/internal/application"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/client/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
}

func NewGrpcController(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g GrpcController) ReadExchanges(ctx context.Context, req *proto.ReadExchangesRequest) (*proto.ReadExchangesResponse, error) {
	list, err := g.application.Commands.CachedReadExchanges.Handle(ctx, nil, req.Names...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.ReadExchangesResponse{
		Exchanges: toGrpcExchanges(list),
	}, nil
}

func toGrpcExchanges(ps []exchange.Exchange) []*proto.Exchange {
	gexchanges := make([]*proto.Exchange, len(ps))
	for i, p := range ps {
		gexchanges[i] = &proto.Exchange{
			Name:         p.Name,
			Pairs:        p.PairsSymbols,
			Periods:      p.PeriodsSymbols,
			Fees:         float32(p.Fees),
			LastSyncTime: p.LastSyncTime.Format(time.RFC3339),
		}
	}
	return gexchanges
}
