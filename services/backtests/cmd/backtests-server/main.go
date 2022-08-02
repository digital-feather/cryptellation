package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers/http/health"
	"github.com/digital-feather/cryptellation/services/backtests/internal/service"
)

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	app, closeApp, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("creating application: %w", err))
		return 255
	}
	defer closeApp()

	// Init grpc server
	grpcController := grpc.New(app)
	if err := grpcController.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running application: %w", err))
		return 255
	}
	defer grpcController.GracefulStop()

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

	log.Print("service is shutting down...")
	return 0
}

func main() {
	os.Exit(run())
}
