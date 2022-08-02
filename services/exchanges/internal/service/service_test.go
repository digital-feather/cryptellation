package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db/cockroach"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/client"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/client/proto"
	"github.com/stretchr/testify/suite"
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
	client    proto.ExchangesServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tmpEnvVar("COCKROACHDB_DATABASE", testDatabase)()
	defer tmpEnvVar("CRYPTELLATION_EXCHANGES_GRPC_URL", ":9003")()

	a, err := newMockApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_EXCHANGES_GRPC_URL")
	grpcController := grpc.New(a)
	suite.NoError(grpcController.RunOnAddr(rpcUrl))

	ok := waitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.Newclient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		go grpcController.Stop() // TODO: remove goroutine
		return closeClient()
	}
}

func (suite *ServiceSuite) SetupTest() {
	defer tmpEnvVar("COCKROACHDB_DATABASE", testDatabase)()

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
	resp, err := suite.client.ReadExchanges(context.Background(), &proto.ReadExchangesRequest{
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
	resp, err = suite.client.ReadExchanges(context.Background(), &proto.ReadExchangesRequest{
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

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}

func waitForPort(address string) bool {
	waitChan := make(chan struct{})

	go func() {
		for {
			conn, err := net.DialTimeout("tcp", address, time.Second)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if conn != nil {
				waitChan <- struct{}{}
				return
			}
		}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-waitChan:
		return true
	case <-timeout:
		return false
	}
}
