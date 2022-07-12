package main

import (
	"fmt"
	"os"

	grpcUtils "github.com/digital-feather/cryptellation/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/controllers"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/service"
	"google.golang.org/grpc"
)

func run() int {
	application, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("creating application: %w", err))
		return 255
	}

	err = grpcUtils.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		exchanges.RegisterExchangesServiceServer(server, svc)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running application: %w", err))
		return 255
	}

	return 0
}

func main() {
	os.Exit(run())
}
