package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcUtils "github.com/digital-feather/cryptellation/internal/go/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/go/controllers/grpc/genproto/exchanges"
	"github.com/digital-feather/cryptellation/internal/go/controllers/http/health"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/controllers"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/service"
	"google.golang.org/grpc"
)

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	app, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("creating application: %w", err))
		return 255
	}

	// Init grpc server
	srv, err := grpcUtils.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(app)
		exchanges.RegisterExchangesServiceServer(server, svc)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running application: %w", err))
		return 255
	}
	defer srv.GracefulStop()

	// Service marked as ready
	log.Println("Service is ready")
	h.Ready(true)

	// Wait for interrupt
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	return 0
}

func main() {
	os.Exit(run())
}
