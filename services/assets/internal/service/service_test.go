package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cryptellation/cryptellation/internal/genproto/assets"
	"github.com/cryptellation/cryptellation/internal/server"
	"github.com/cryptellation/cryptellation/internal/tests"
	"github.com/cryptellation/cryptellation/pkg/client"
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db/cockroach"
	"github.com/cryptellation/cryptellation/services/assets/internal/application"
	"github.com/cryptellation/cryptellation/services/assets/internal/controllers"
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
	client    assets.AssetsServiceClient
	closeTest func()
}

func (suite *ServiceSuite) BeforeTest(suiteName, testName string) {
	defer tests.TempEnvVar("COCKROACHDB_DATABASE", "assets")()
	defer tests.TempEnvVar("CRYPTELLATION_ASSETS_GRPC_URL", ":9000")()

	a, closeApp, err := NewApplication()
	suite.Require().NoError(err)
	suite.app = a

	rpcUrl := os.Getenv("CRYPTELLATION_ASSETS_GRPC_URL")
	go server.RunGRPCServerOnAddr(rpcUrl, func(server *grpc.Server) {
		svc := controllers.NewGrpcController(a)
		assets.RegisterAssetsServiceServer(server, svc)
	})

	ok := tests.WaitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.NewAssetsGrpcClient()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() {
		closeClient()
		closeApp()
	}

	suite.Require().NoError(cockroach.Reset())
}

func (suite *ServiceSuite) AfterTest(suiteName, testName string) {
	suite.closeTest()
}

func (suite *ServiceSuite) TestCreateReadAssets() {
	_, err := suite.client.CreateAssets(context.Background(), &assets.CreateAssetsRequest{
		Assets: []*assets.Asset{
			{Symbol: "ETH"},
			{Symbol: "BTC"},
		},
	})
	suite.Require().NoError(err)

	resp, err := suite.client.ReadAssets(context.Background(), &assets.ReadAssetsRequest{
		Symbols: []string{
			"ETH",
			"BTC",
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Assets, 2)
	suite.Require().Contains(resp.Assets, &assets.Asset{
		Symbol: "ETH",
	})
	suite.Require().Contains(resp.Assets, &assets.Asset{
		Symbol: "BTC",
	})
}
