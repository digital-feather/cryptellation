package redis

import (
	"errors"
	"os"

	"golang.org/x/xerrors"
)

var (
	ErrInvalidConfig = errors.New("invalid redis config")
)

type Config struct {
	Address  string
	Password string
}

func (c *Config) Load() *Config {
	c.Address = os.Getenv("REDIS_ADDRESS")
	c.Password = os.Getenv("REDIS_PASSWORD")

	return c
}

func (c Config) Validate() error {
	if c.Address == "" {
		return xerrors.Errorf("reading address from env (%q): %w", c.Address, ErrInvalidConfig)
	}

	return nil
}
