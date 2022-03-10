package service

import (
	"log"
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/pkg/client"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/db/cockroach"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/application"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/controllers"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

func TestServiceSuite(t *testing.T) {
	if os.Getenv("COCKROACHDB_HOST") == "" {
		t.Skip()
	}

	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	app       application.Application
	client    candlesticks.CandlesticksServiceClient
	closeTest func()
}

func (suite *ServiceSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "candlesticks")()
	defer tests.TempEnvVar("CRYPTELLATION_CANDLESTICKS_GRPC_URL", ":9002")()

	a, err := newMockApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_CANDLESTICKS_GRPC_URL")
	go server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		candlesticks.RegisterCandlesticksServiceServer(server, svc)
	})

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewCandlesticksGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() {
		closeClient()
	}

	suite.Require().NoError(cockroach.Reset())
}

func (suite *ServiceSuite) AfterTest(suiteName, testName string) {
	suite.closeTest()
}
