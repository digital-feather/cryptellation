package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	grpcUtils "github.com/digital-feather/cryptellation/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/ticks"
	"github.com/digital-feather/cryptellation/internal/tests"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

func TestServiceSuite(t *testing.T) {
	if os.Getenv("REDIS_ADDRESS") == "" {
		t.Skip()
	}

	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	vdb       vdb.Port
	client    ticks.TicksServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupTest() {
	defer tests.TempEnvVar("CRYPTELLATION_TICKS_GRPC_URL", ":9005")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_TICKS_GRPC_URL")
	grpcServer, err := grpcUtils.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		ticks.RegisterTicksServiceServer(server, svc)
	})
	suite.Require().NoError(err)

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewTicksGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		err = closeClient()
		go grpcServer.Stop() // TODO: remove goroutine
		closeApplication()
		return err
	}

	vdb, err := redis.New()
	suite.Require().NoError(err)
	suite.vdb = vdb
}

func (suite *ServiceSuite) AfterTest(suiteName, testName string) {
	err := suite.closeTest()
	suite.NoError(err)
}

func (suite *ServiceSuite) TestListenSymbol() {
	stream, err := suite.client.ListenSymbol(context.Background(),
		&ticks.ListenSymbolRequest{
			Exchange:   "mock_exchange",
			PairSymbol: "SYMBOL",
		})
	suite.Require().NoError(err)

	for i := int64(0); i < 50; i++ {
		t, err := stream.Recv()
		suite.Require().NoError(err)
		suite.Require().Equal("mock_exchange", t.Exchange)
		suite.Require().Equal("SYMBOL", t.PairSymbol)
		suite.Require().Equal(float32(100+i), t.Price)
		ti, err := time.Parse(time.RFC3339Nano, t.Time)
		suite.Require().NoError(err)
		suite.Require().WithinDuration(time.Unix(i, 0), ti, time.Microsecond)
	}
}
