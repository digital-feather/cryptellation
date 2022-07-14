package main

import (
	"fmt"
	"os"

	grpcUtils "github.com/digital-feather/cryptellation/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/backtests"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers"
	"github.com/digital-feather/cryptellation/services/backtests/internal/service"
	"google.golang.org/grpc"
)

func run() int {
	application, closeFunc, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("creating application: %w", err))
		return 255
	}
	defer func() {
		if err := closeFunc(); err != nil {
			fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("closing application: %w", err))
		}
	}()

	err = grpcUtils.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		backtests.RegisterBacktestsServiceServer(server, svc)
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
