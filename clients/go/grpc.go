package client

import (
	"os"

	"github.com/digital-feather/cryptellation/internal/genproto/backtests"
	"github.com/digital-feather/cryptellation/internal/genproto/candlesticks"
	"github.com/digital-feather/cryptellation/internal/genproto/exchanges"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func NewExchangesGrpcClient() (client exchanges.ExchangesServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, xerrors.Errorf("dialing exchanges grpc server: %w", err)
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
		return nil, func() error { return nil }, xerrors.Errorf("dialing candlesticks grpc server: %w", err)
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
		return nil, func() error { return nil }, xerrors.Errorf("dialing backtests grpc server: %w", err)
	}

	return backtests.NewBacktestsServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithInsecure()}
}
