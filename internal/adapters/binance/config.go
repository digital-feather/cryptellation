package binance

import (
	"errors"
	"os"

	"golang.org/x/xerrors"
)

var (
	ErrInvalidConfig = errors.New("invalid binance config")
)

type Config struct {
	ApiKey    string
	SecretKey string
}

func (c *Config) Load() *Config {
	c.ApiKey = os.Getenv("BINANCE_API_KEY")
	c.SecretKey = os.Getenv("BINANCE_SECRET_KEY")

	return c
}

func (c Config) Validate() error {
	if c.ApiKey == "" {
		return xerrors.Errorf("reading api key from env (%q): %w", c.ApiKey, ErrInvalidConfig)
	}

	if c.SecretKey == "" {
		return xerrors.Errorf("reading secret key from env (%q): %w", c.SecretKey, ErrInvalidConfig)
	}

	return nil
}
