package redis

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
		Address, Password string
		Err               error
	}{
		{
			Address:  "address",
			Password: "password",
		},
	}

	var config Config
	for i, c := range cases {
		setEnv(c.Address, c.Password)

		err := config.Load().Validate()
		suite.Require().Equal(c.Err, err, i)

		setEnv("", "")
	}
}

func setEnv(address, password string) {
	os.Setenv("REDIS_ADDRESS", address)
	os.Setenv("REDIS_PASSWORD", password)
}
