package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		panic("no port")
	}
	addr := fmt.Sprintf(":%s", port)
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	grpcServer := grpc.NewServer()
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Starting: gRPC Listener")
	log.Fatal(grpcServer.Serve(listen))
}
