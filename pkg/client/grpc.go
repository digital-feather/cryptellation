package client

import (
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
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
		return nil, func() error { return nil }, xerrors.Errorf("dialing asset grpc server: %w", err)
	}

	return assets.NewAssetsServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(grpcAddr string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithInsecure()}
}
