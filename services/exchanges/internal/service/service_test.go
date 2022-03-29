package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	client "github.com/cryptellation/cryptellation/clients/go"
	"github.com/cryptellation/cryptellation/internal/genproto/exchanges"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/adapters/db/cockroach"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/application"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/controllers"
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
	client    exchanges.ExchangesServiceClient
	closeTest func()
}

func (suite *ServiceSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "exchanges")()
	defer tests.TempEnvVar("CRYPTELLATION_EXCHANGES_GRPC_URL", ":9003")()

	a, err := newMockApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	go server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		exchanges.RegisterExchangesServiceServer(server, svc)
	})

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewExchangesGrpcClient()
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

func (suite *ServiceSuite) TestReadExchanges() {
	// When requesting an exchange for the first time
	resp, err := suite.client.ReadExchanges(context.Background(), &exchanges.ReadExchangesRequest{
		Names: []string{
			"mock_exchange",
		},
	})

	// Then the request is successful
	suite.Require().NoError(err)

	// And the exchange is correct
	suite.Require().Len(resp.Exchanges, 1)
	suite.Require().Equal("mock_exchange", resp.Exchanges[0].Name)

	// And the last sync time is now
	firstTime := resp.Exchanges[0].LastSyncTime
	t, err := time.Parse(time.RFC3339, firstTime)
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Now().UTC(), t, 2*time.Second)

	// When the second request is made
	resp, err = suite.client.ReadExchanges(context.Background(), &exchanges.ReadExchangesRequest{
		Names: []string{
			"mock_exchange",
		},
	})

	// Then the request is successful
	suite.Require().NoError(err)

	// And the last sync time is the same as previous
	suite.Require().Len(resp.Exchanges, 1)
	suite.Require().Equal(firstTime, resp.Exchanges[0].LastSyncTime)
}
