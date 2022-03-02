package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/genproto/pairs"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/pkg/client"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db/cockroach"
	"github.com/cryptellation/cryptellation/services/pairs/internal/application"
	"github.com/cryptellation/cryptellation/services/pairs/internal/controllers"
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
	client    pairs.PairsServiceClient
	closeTest func()
}

func (suite *ServiceSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "pairs")()
	defer tests.TempEnvVar("CRYPTELLATION_PAIRS_GRPC_URL", ":9001")()

	a, err := NewApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_PAIRS_GRPC_URL")
	go server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		pairs.RegisterPairsServiceServer(server, svc)
	})

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewPairsGrpcClient()
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

func (suite *ServiceSuite) TestCreateReadPairs() {
	_, err := suite.client.CreatePairs(context.Background(), &pairs.CreatePairsRequest{
		Pairs: []*pairs.Pair{
			{
				Symbol:           "ETH-USDC",
				BaseAssetSymbol:  "ETH",
				QuoteAssetSymbol: "USDC",
			},
		},
	})
	suite.Require().NoError(err)

	resp, err := suite.client.ReadPairs(context.Background(), &pairs.ReadPairsRequest{
		Symbols: []string{
			"ETH-USDC",
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Pairs, 1)
	suite.Require().Contains(resp.Pairs, &pairs.Pair{
		Symbol:           "ETH-USDC",
		BaseAssetSymbol:  "ETH",
		QuoteAssetSymbol: "USDC",
	})
}
