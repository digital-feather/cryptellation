package service

import (
	"context"
	"log"
	"os"
	"testing"

	client "github.com/digital-feather/cryptellation/clients/go"
	grpcUtils "github.com/digital-feather/cryptellation/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/controllers/grpc/genproto/livetests"
	"github.com/digital-feather/cryptellation/internal/tests"
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/livetests/internal/controllers"
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
	client    livetests.LivetestsServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tests.TempEnvVar("CRYPTELLATION_LIVETESTS_GRPC_URL", ":9006")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_LIVETESTS_GRPC_URL")
	grpcServer, err := grpcUtils.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		livetests.RegisterLivetestsServiceServer(server, svc)
	})
	suite.NoError(err)

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewLivetestsGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		err := closeClient()
		go grpcServer.Stop() // TODO: remove goroutine
		closeApplication()
		return err
	}

	vdb, err := redis.New()
	suite.Require().NoError(err)
	suite.vdb = vdb
}

func (suite *ServiceSuite) TearDownSuite() {
	err := suite.closeTest()
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) TestCreateLivetest() {
	req := livetests.CreateLivetestRequest{
		Accounts: map[string]*livetests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateLivetest(context.Background(), &req)
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadLivetest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.Accounts, 1)
	suite.Require().Len(recvBT.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(1000.0, recvBT.Accounts["exchange"].Balances["DAI"])
}

func (suite *ServiceSuite) TestLivetestSubscribeToEvents() {
	req := livetests.CreateLivetestRequest{
		Accounts: map[string]*livetests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateLivetest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToLivetestEvents(context.Background(), &livetests.SubscribeToLivetestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadLivetest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.TickSubscribers, 1)
	suite.Require().Equal("exchange", recvBT.TickSubscribers[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", recvBT.TickSubscribers[0].PairSymbol)
}
