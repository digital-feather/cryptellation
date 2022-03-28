package service

import (
	"context"
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
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/event"
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

	a, closeApplication, err := NewMockedApplication()
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

	client, closeClient, err := client.NewBacktestsGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() {
		closeClient()
		closeApplication()
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

func (suite *ServiceSuite) TestBacktestSubscribeToEvents() {
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

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &backtests.SubscribeToBacktestEventsRequest{
		Id:         resp.Id,
		Exchange:   "exchange",
		PairSymbol: "ETH-DAI",
	})
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(resp.Id))
	suite.Require().Len(recvBT.TickSubscribers, 1)
	suite.Require().Equal("exchange", recvBT.TickSubscribers[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", recvBT.TickSubscribers[0].PairSymbol)
}

func (suite *ServiceSuite) TestBacktestListenEvents() {
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

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &backtests.SubscribeToBacktestEventsRequest{
		Id:         resp.Id,
		Exchange:   "exchange",
		PairSymbol: "ETH-DAI",
	})
	suite.Require().NoError(err)

	stream, err := suite.client.ListenBacktest(context.Background(), &backtests.ListenBacktestRequest{
		Id: resp.Id,
	})
	suite.Require().NoError(err)

	// First candlestick
	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:00:00Z", "null")

	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":2,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:00:00Z", "null")

	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":0.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:00:00Z", "null")

	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:00:00Z", "null")

	// Second candlestick
	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:01:00Z", "null")

	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":2,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:01:00Z", "null")

	suite.advance(resp.Id, false)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":0.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:01:00Z", "null")

	suite.advance(resp.Id, true)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsEnd, "1970-01-01T00:01:00Z", "null")
}

func (suite *ServiceSuite) advance(id uint64, finished bool) {
	resp, err := suite.client.AdvanceBacktest(context.Background(), &backtests.AdvanceBacktestRequest{
		Id: id,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(finished, resp.Finished)
}

func (suite *ServiceSuite) checkEvent(stream backtests.BacktestsService_ListenBacktestClient, evtType event.Type, time, content string) {
	evt, err := stream.Recv()
	suite.Require().NoError(err)
	suite.Require().Equal(evtType.String(), evt.Type)
	suite.Require().Equal(time, evt.Time)
	suite.Require().Equal(content, evt.Content)
}
