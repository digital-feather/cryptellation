package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	grpcUtils "github.com/digital-feather/cryptellation/internal/go/controllers/grpc"
	"github.com/digital-feather/cryptellation/internal/go/controllers/grpc/genproto/exchanges"
	"github.com/digital-feather/cryptellation/internal/go/tests"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db/cockroach"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/controllers"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const (
	testDatabase = "exchanges"
)

func TestServiceSuite(t *testing.T) {
	if os.Getenv("COCKROACHDB_HOST") == "" {
		t.Skip()
	}

	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	db        *cockroach.DB
	client    exchanges.ExchangesServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", testDatabase)()
	defer tests.TempEnvVar("CRYPTELLATION_EXCHANGES_GRPC_URL", ":9003")()

	a, err := newMockApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	grpcServer, err := grpcUtils.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		exchanges.RegisterExchangesServiceServer(server, svc)
	})
	suite.Require().NoError(err)

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewExchangesGrpcClient()
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
