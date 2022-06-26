package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/genproto/backtests"
	"github.com/digital-feather/cryptellation/internal/server"
	"github.com/digital-feather/cryptellation/internal/tests"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
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
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tests.TempEnvVar("CRYPTELLATION_BACKTESTS_GRPC_URL", ":9004")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_BACKTESTS_GRPC_URL")
	go func() {
		err := server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
			svc := controllers.NewGrpcController(a)
			backtests.RegisterBacktestsServiceServer(server, svc)
		})
		suite.NoError(err)
	}()

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewBacktestsGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		err := closeClient()
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

func (suite *ServiceSuite) TestCreateBacktest() {
	req := backtests.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: map[string]*backtests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
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
		Accounts: map[string]*backtests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &backtests.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.TickSubscribers, 1)
	suite.Require().Equal("exchange", recvBT.TickSubscribers[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", recvBT.TickSubscribers[0].PairSymbol)
}

func (suite *ServiceSuite) TestBacktestListenEvents() {
	req := backtests.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: map[string]*backtests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &backtests.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	stream, err := suite.client.ListenBacktest(context.Background())
	suite.Require().NoError(err)

	// First candlestick (high)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":2,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:00:00Z", "{\"finished\":false}")

	// First candlestick (low)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":0.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:00:00Z", "{\"finished\":false}")

	// First candlestick (close)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:00:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:00:00Z", "{\"finished\":false}")

	// Second candlestick (open)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:01:00Z", "{\"finished\":false}")

	// Second candlestick (high)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":2,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:01:00Z", "{\"finished\":false}")

	// Second candlestick (low)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":0.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:01:00Z", "{\"finished\":false}")

	// Second candlestick (close)
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsTick, "1970-01-01T00:01:00Z", "{\"pair_symbol\":\"ETH-DAI\",\"price\":1.5,\"exchange\":\"exchange\"}")
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:01:00Z", "{\"finished\":false}")

	// End of backtest
	suite.advance(stream, resp.Id)
	suite.checkEvent(stream, event.TypeIsStatus, "1970-01-01T00:02:00Z", "{\"finished\":true}")
}

func (suite *ServiceSuite) advance(stream backtests.BacktestsService_ListenBacktestClient, id uint64) {
	err := stream.Send(&backtests.BacktestEventRequest{
		Id: id,
	})
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) checkEvent(stream backtests.BacktestsService_ListenBacktestClient, evtType event.Type, time, content string) {
	evt, err := stream.Recv()
	suite.Require().NoError(err)
	suite.Require().Equal(evtType.String(), evt.Type)
	suite.Require().Equal(time, evt.Time)
	suite.Require().Equal(content, evt.Content)
}

func (suite *ServiceSuite) passEvent(stream backtests.BacktestsService_ListenBacktestClient, evtType event.Type) {
	evt, err := stream.Recv()
	suite.Require().NoError(err)
	suite.Require().Equal(evtType.String(), evt.Type)
}

func (suite *ServiceSuite) TestBacktestOrders() {
	req := backtests.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(600, 0).Format(time.RFC3339),
		Accounts: map[string]*backtests.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &backtests.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	_, err = suite.client.CreateBacktestOrder(context.Background(), &backtests.CreateBacktestOrderRequest{
		BacktestId:   resp.Id,
		Type:         "market",
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         "buy",
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accountsResp, err := suite.client.Accounts(context.Background(), &backtests.AccountsRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(float32(999), accountsResp.Accounts["exchange"].Assets["DAI"])
	suite.Require().Equal(float32(1), accountsResp.Accounts["exchange"].Assets["ETH"])

	stream, err := suite.client.ListenBacktest(context.Background())
	suite.Require().NoError(err)
	for i := 0; i < 5; i++ {
		suite.advance(stream, resp.Id)
		suite.passEvent(stream, event.TypeIsTick)
		suite.passEvent(stream, event.TypeIsStatus)
	}

	_, err = suite.client.CreateBacktestOrder(context.Background(), &backtests.CreateBacktestOrderRequest{
		BacktestId:   resp.Id,
		Type:         "market",
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         "sell",
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accountsResp, err = suite.client.Accounts(context.Background(), &backtests.AccountsRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(float32(1001), accountsResp.Accounts["exchange"].Assets["DAI"])
	suite.Require().Equal(float32(0), accountsResp.Accounts["exchange"].Assets["ETH"])

	ordersResp, err := suite.client.Orders(context.Background(), &backtests.OrdersRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Len(ordersResp.Orders, 2)

	suite.Require().Equal("1970-01-01T00:00:00Z", ordersResp.Orders[0].Time)
	suite.Require().Equal("market", ordersResp.Orders[0].Type)
	suite.Require().Equal("exchange", ordersResp.Orders[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", ordersResp.Orders[0].PairSymbol)
	suite.Require().Equal("buy", ordersResp.Orders[0].Side)
	suite.Require().Equal(float32(1), ordersResp.Orders[0].Quantity)
	suite.Require().Equal(float32(1), ordersResp.Orders[0].Price)

	suite.Require().Equal("1970-01-01T00:01:00Z", ordersResp.Orders[1].Time)
	suite.Require().Equal("market", ordersResp.Orders[1].Type)
	suite.Require().Equal("exchange", ordersResp.Orders[1].ExchangeName)
	suite.Require().Equal("ETH-DAI", ordersResp.Orders[1].PairSymbol)
	suite.Require().Equal("sell", ordersResp.Orders[1].Side)
	suite.Require().Equal(float32(1), ordersResp.Orders[1].Quantity)
	suite.Require().Equal(float32(2), ordersResp.Orders[1].Price)
}
