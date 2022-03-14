package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	client "github.com/cryptellation/cryptellation/clients/go"
	"github.com/cryptellation/cryptellation/internal/genproto/backtests"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb/redis"
	"github.com/cryptellation/cryptellation/services/backtests/internal/application"
	"github.com/cryptellation/cryptellation/services/backtests/internal/controllers"
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
	app       application.Application
	vdb       vdb.Port
	client    backtests.BacktestsServiceClient
	closeTest func()
}

func (suite *ServiceSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("CRYPTELLATION_BACKTESTS_GRPC_URL", ":9004")()

	a, err := NewApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_BACKTESTS_GRPC_URL")
	go server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		backtests.RegisterBacktestsServiceServer(server, svc)
	})

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}
	fmt.Println("ok")

	client, closeClient, err := client.NewBacktestsGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() {
		closeClient()
	}

	vdb, err := redis.New()
	suite.Require().NoError(err)
	suite.vdb = vdb
}

func (suite *ServiceSuite) AfterTest(suiteName, testName string) {
	suite.closeTest()
}

func (suite *ServiceSuite) TestCreateBacktest() {
	req := backtests.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: []*backtests.Account{
			{
				ExchangeName: "exchange",
				Assets: []*backtests.AssetQuantity{
					{
						AssetName: "DAI",
						Quantity:  1000,
					},
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(0, 0), recvBT.StartTime, time.Millisecond)
	suite.Require().WithinDuration(time.Unix(120, 0), recvBT.EndTime, time.Millisecond)
	suite.Require().Len(recvBT.Accounts, 1)
	suite.Require().Len(recvBT.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(1000.0, recvBT.Accounts["exchange"].Balances["DAI"])
}
