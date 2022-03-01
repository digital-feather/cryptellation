package client

import (
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
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

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithInsecure()}
}