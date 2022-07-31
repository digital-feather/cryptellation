package binance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) TestLoadValidate() {
	cases := []struct {
		Api, Secret string
		Err         error
	}{
		{
			Api:    "api-key",
			Secret: "secret-key",
		},
	}

	var config Config
	for i, c := range cases {
		setEnv(c.Api, c.Secret)

		err := config.Load().Validate()
		suite.Require().Equal(c.Err, err, i)

		setEnv("", "")
	}
}

func setEnv(api, secret string) {
	os.Setenv("BINANCE_API_KEY", api)
	os.Setenv("BINANCE_SECRET_KEY", secret)
}
