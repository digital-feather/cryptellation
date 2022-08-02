package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	app "github.com/digital-feather/cryptellation/services/livetests/internal/application"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
	"github.com/digital-feather/cryptellation/services/livetests/pkg/client/proto"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
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
	proto.RegisterLivetestsServiceServer(grpcServer, g)

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

func (g GrpcController) CreateLivetest(ctx context.Context, req *proto.CreateLivetestRequest) (*proto.CreateLivetestResponse, error) {
	newPayload, err := fromCreateLivetestRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := g.application.Commands.Livetest.Create.Handle(ctx, newPayload)
	if err != nil {
		return nil, err
	}

	return &proto.CreateLivetestResponse{
		Id: uint64(id),
	}, nil
}

func fromCreateLivetestRequest(req *proto.CreateLivetestRequest) (livetest.NewPayload, error) {
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

func (g GrpcController) SubscribeToLivetestEvents(ctx context.Context, req *proto.SubscribeToLivetestEventsRequest) (*proto.SubscribeToLivetestEventsResponse, error) {
	err := g.application.Commands.Livetest.SubscribeToEvents.Handle(ctx, uint(req.Id), req.ExchangeName, req.PairSymbol)
	return &proto.SubscribeToLivetestEventsResponse{}, err
}

func (g GrpcController) ListenLivetest(srv proto.LivetestsService_ListenLivetestServer) error {
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

func sendEventResponse(srv proto.LivetestsService_ListenLivetestServer, evt event.Event) error {
	content, err := json.Marshal(evt.Content)
	if err != nil {
		return fmt.Errorf("marshaling event content: %w", err)
	}

	return srv.Send(&proto.LivetestEventResponse{
		Type:    evt.Type.String(),
		Time:    evt.Time.Format(time.RFC3339),
		Content: string(content),
	})
}
