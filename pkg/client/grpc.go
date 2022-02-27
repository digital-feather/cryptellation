package client

import (
	"errors"
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"google.golang.org/grpc"
)

func NewAssetsGrpcClient() (client assets.AssetsServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("CRYPTELLATION_ASSETS_GRPC_URL")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("CRYPTELLATION_ASSETS_GRPC_URL")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return assets.NewAssetsServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(grpcAddr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{grpc.WithInsecure()}, nil
}
