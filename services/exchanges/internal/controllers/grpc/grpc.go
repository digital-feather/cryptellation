package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	app "github.com/digital-feather/cryptellation/services/exchanges/internal/application"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/client/proto"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcController struct {
	application app.Application
	server      *grpc.Server
}

func New(application app.Application) GrpcController {
	return GrpcController{application: application}
}

func (g *GrpcController) Run() error {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		return xerrors.New("no service port provided")
	}
	addr := fmt.Sprintf(":%s", port)
	return g.RunOnAddr(addr)
}

func (g *GrpcController) RunOnAddr(addr string) error {
	grpcServer := grpc.NewServer()
	proto.RegisterExchangesServiceServer(grpcServer, g)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listening error: %w", err)
	}

	log.Println("Starting: gRPC Listener")
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Println("error when serving grpc:", err)
		}
	}()

	return nil
}

func (g *GrpcController) GracefulStop() {
	if g.server == nil {
		log.Println("WARNING: attempted to gracefully stop a non running grpc server")
		return
	}

	g.server.GracefulStop()
	g.server = nil
}

func (g *GrpcController) Stop() {
	if g.server == nil {
		log.Println("WARNING: attempted to stop a non running grpc server")
		return
	}

	g.server.Stop()
	g.server = nil
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
