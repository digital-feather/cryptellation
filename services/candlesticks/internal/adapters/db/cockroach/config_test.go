package cockroach

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
		Host, Port, User, Password, Database string
		Err                                  error
	}{
		{
			Host:     "host",
			Port:     "1000",
			User:     "user",
			Password: "password",
			Database: "database",
			Err:      nil,
		},
	}

	var config Config
	for i, c := range cases {
		setEnv(c.Host, c.Port, c.User, c.Password, c.Database)

		err := config.Load().Validate()
		suite.Require().Equal(c.Err, err, i)

		setEnv("", "", "", "", "")
	}
}

func setEnv(host, port, user, password, database string) {
	os.Setenv("COCKROACHDB_HOST", host)
	os.Setenv("COCKROACHDB_PORT", port)
	os.Setenv("COCKROACHDB_USER", user)
	os.Setenv("COCKROACHDB_PASSWORD", password)
	os.Setenv("COCKROACHDB_DATABASE", database)
}
