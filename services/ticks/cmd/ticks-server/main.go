package main

import (
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/internal/genproto/ticks"
	"github.com/digital-feather/cryptellation/internal/server"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers"
	"github.com/digital-feather/cryptellation/services/ticks/internal/service"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func run() int {
	application, closeFunc, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", xerrors.Errorf("creating application: %w", err))
		return 255
	}
	defer closeFunc()

	server.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		ticks.RegisterTicksServiceServer(server, svc)
	})

	return 0
}

func main() {
	os.Exit(run())
}
