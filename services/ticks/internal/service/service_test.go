package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
	"github.com/stretchr/testify/suite"
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
	client    proto.TicksServiceClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupTest() {
	defer tmpEnvVar("CRYPTELLATION_TICKS_GRPC_URL", ":9005")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_TICKS_GRPC_URL")
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
		err = closeClient()
		go grpcController.Stop() // TODO: remove goroutine
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
		&proto.ListenSymbolRequest{
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
