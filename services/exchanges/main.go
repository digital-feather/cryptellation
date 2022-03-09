package main

import (
	"fmt"
	"os"

	"github.com/cryptellation/cryptellation/internal/genproto/exchanges"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/controllers"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/service"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func run() int {
	application, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", xerrors.Errorf("creating application: %w", err))
		return 255
	}

	server.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		exchanges.RegisterExchangesServiceServer(server, svc)
	})

	return 0
}

func main() {
	os.Exit(run())
}
