package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	app "github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/commands"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
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
	proto.RegisterCandlesticksServiceServer(grpcServer, g)

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

func (g GrpcController) ReadCandlesticks(ctx context.Context, req *proto.ReadCandlesticksRequest) (*proto.ReadCandlesticksResponse, error) {
	payload, err := fromReadCandlesticksRequest(req)
	if err != nil {
		return nil, err
	}

	list, err := g.application.Commands.CachedReadCandlesticks.Handle(ctx, payload)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.ReadCandlesticksResponse{
		Candlesticks: toGrpcCandlesticks(list),
	}, nil
}

func fromReadCandlesticksRequest(req *proto.ReadCandlesticksRequest) (commands.CachedReadCandlesticksPayload, error) {
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

func toGrpcCandlesticks(cl *candlestick.List) []*proto.Candlestick {
	gcandlesticks := make([]*proto.Candlestick, 0, cl.Len())
	_ = cl.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		gcandlesticks = append(gcandlesticks, &proto.Candlestick{
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
