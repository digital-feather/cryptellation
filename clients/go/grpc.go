package client

import (
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/backtests"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/candlesticks"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/exchanges"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/ticks"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExchangesGrpcClient() (client exchanges.ExchangesServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing exchanges grpc server: %w", err)
	}

	return exchanges.NewExchangesServiceClient(conn), conn.Close, nil
}

func NewCandlesticksGrpcClient() (client candlesticks.CandlesticksServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_CANDLESTICKS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing candlesticks grpc server: %w", err)
	}

	return candlesticks.NewCandlesticksServiceClient(conn), conn.Close, nil
}

func NewBacktestsGrpcClient() (client backtests.BacktestsServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_BACKTESTS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing backtests grpc server: %w", err)
	}

	return backtests.NewBacktestsServiceClient(conn), conn.Close, nil
}

func NewTicksGrpcClient() (client ticks.TicksServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_TICKS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("dialing ticks grpc server: %w", err)
	}

	return ticks.NewTicksServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
