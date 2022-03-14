package client

import (
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/internal/genproto/backtests"
	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	"github.com/cryptellation/cryptellation/internal/genproto/exchanges"
	"github.com/cryptellation/cryptellation/internal/genproto/pairs"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func NewAssetsGrpcClient() (client assets.AssetsServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_ASSETS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, xerrors.Errorf("dialing assets grpc server: %w", err)
	}

	return assets.NewAssetsServiceClient(conn), conn.Close, nil
}

func NewPairsGrpcClient() (client pairs.PairsServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_PAIRS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, xerrors.New("no grpc url provided")
	}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts(grpcAddr)...)
	if err != nil {
		return nil, func() error { return nil }, xerrors.Errorf("dialing pairs grpc server: %w", err)
	}

	return pairs.NewPairsServiceClient(conn), conn.Close, nil
}

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
