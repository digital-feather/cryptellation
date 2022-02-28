package main

import (
	"fmt"
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/services/assets/internal/controllers"
	"github.com/cryptellation/cryptellation/services/assets/internal/service"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func run() int {
	application, cleanup, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", xerrors.Errorf("creating application: %w", err))
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
