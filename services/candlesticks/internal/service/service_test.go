package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	grpcUtils "github.com/digital-feather/cryptellation/internal/go/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/go/tests"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db/cockroach"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const (
	testDatabase = "candlesticks"
)

func TestServiceSuite(t *testing.T) {
	if os.Getenv("COCKROACHDB_HOST") == "" {
		t.Skip()
	}

	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	db        db.Port
	client    proto.CandlesticksServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", testDatabase)()
	defer tests.TempEnvVar("CRYPTELLATION_CANDLESTICKS_GRPC_URL", ":9002")()

	a, err := newMockApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_CANDLESTICKS_GRPC_URL")
	grpcServer, err := grpcUtils.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		proto.RegisterCandlesticksServiceServer(server, svc)
	})
	suite.Require().NoError(err)

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.Newclient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		go grpcServer.Stop() // TODO: remove goroutine
		return closeClient()
	}
}

func (suite *ServiceSuite) SetupTest() {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", testDatabase)()

	db, err := cockroach.New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.db = db
}

func (suite *ServiceSuite) TearDownSuite() {
	err := suite.closeTest()
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) TestGetCandlesticksAllExistWithNoneInDB() {
	// Given a client service
	// Provided before

	// When a request is made
	resp, err := suite.client.ReadCandlesticks(context.Background(), &proto.ReadCandlesticksRequest{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		PeriodSymbol: period.M1.String(),
		Start:        time.Unix(0, 0).Format(time.RFC3339),
		End:          time.Unix(540, 0).Format(time.RFC3339),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Len(resp.Candlesticks, 10)
	for i, cs := range resp.Candlesticks {
		suite.Require().Equal(float32(60*i), cs.Open)
		suite.Require().Equal(time.Unix(int64(60*i), 0).Format(time.RFC3339Nano), cs.Time)
	}
}

func (suite *ServiceSuite) TestGetCandlesticksAllInexistantWithNoneInDB() {
	// Given a client service
	// Provided before

	// When a request is made
	resp, err := suite.client.ReadCandlesticks(context.Background(), &proto.ReadCandlesticksRequest{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		PeriodSymbol: period.M1.String(),
		Start:        time.Unix(60000, 0).Format(time.RFC3339),
		End:          time.Unix(60600, 0).Format(time.RFC3339),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Len(resp.Candlesticks, 0)
}

func (suite *ServiceSuite) TestGetCandlesticksFromDBAndService() {
	// Given a client service
	// Provided before

	// And candlesticks in DB
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	for i := int64(0); i < 10; i++ {
		err := cl.Set(time.Unix(60*i, 0), candlestick.Candlestick{
			Open:  float64(i),
			Close: 4321,
		})
		suite.Require().NoError(err)
	}
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), cl))

	// When a request is made
	resp, err := suite.client.ReadCandlesticks(context.Background(), &proto.ReadCandlesticksRequest{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		PeriodSymbol: period.M1.String(),
		Start:        time.Unix(0, 0).Format(time.RFC3339),
		End:          time.Unix(1140, 0).Format(time.RFC3339),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Len(resp.Candlesticks, 20)
	for i, cs := range resp.Candlesticks {
		suite.Require().Equal(time.Unix(int64(60*i), 0).Format(time.RFC3339Nano), cs.Time)
		if i < 10 {
			suite.Require().Equal(float32(4321), cs.Close, i)
		} else {
			suite.Require().Equal(float32(1234), cs.Close, i)
		}
	}
}

func (suite *ServiceSuite) TestGetCandlesticksFromDBAndServiceWithUncomplete() {
	// Given a client service
	// Provided before

	// And candlesticks in DB
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	for i := int64(0); i < 10; i++ {
		cs := candlestick.Candlestick{
			Open:  float64(i),
			Close: 4321,
		}

		if i == 5 {
			cs.Uncomplete = true
		}

		err := cl.Set(time.Unix(60*i, 0), cs)
		suite.Require().NoError(err)
	}
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), cl))

	// When a request is made
	resp, err := suite.client.ReadCandlesticks(context.Background(), &proto.ReadCandlesticksRequest{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		PeriodSymbol: period.M1.String(),
		Start:        time.Unix(0, 0).Format(time.RFC3339),
		End:          time.Unix(1140, 0).Format(time.RFC3339),
	})

	// Then all candlesticks are retrieved from mock
	suite.Require().NoError(err)
	suite.Require().Len(resp.Candlesticks, 20)
	for i, cs := range resp.Candlesticks {
		suite.Require().Equal(time.Unix(int64(60*i), 0).Format(time.RFC3339Nano), cs.Time)
		suite.Require().Equal(float32(1234), cs.Close, i)
	}
}
