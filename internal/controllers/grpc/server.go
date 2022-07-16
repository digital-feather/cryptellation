package grpc

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) (*grpc.Server, error) {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		return nil, xerrors.New("no service port provided")
	}
	addr := fmt.Sprintf(":%s", port)
	return RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("grpc listening error: %w", err)
	}

	log.Println("Starting: gRPC Listener")
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Println("error when serving grpc:", err)
		}
	}()

	return grpcServer, nil
}
