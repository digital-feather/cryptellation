package main

import (
	"fmt"
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/services/assets/internal/controllers"
	"github.com/cryptellation/cryptellation/services/assets/internal/service"
	"google.golang.org/grpc"
)

func run() int {
	application, err, cleanup := service.NewApplication()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 255
	}
	defer cleanup()

	server.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		assets.RegisterAssetsServiceServer(server, svc)
	})

	return 0
}

func main() {
	os.Exit(run())
}
