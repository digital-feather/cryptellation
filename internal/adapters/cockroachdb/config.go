package cockroachdb

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DefaultGormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
)

var (
	ErrNoUserInConfig = errors.New("no-cockroach-user-in-config")
	ErrNoHostInConfig = errors.New("no-cockroach-host-in-config")
	ErrNoPortInConfig = errors.New("no-cockroach-port-in-config")
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func (c *Config) Load() *Config {
	c.Host = os.Getenv("COCKROACHDB_HOST")
	c.Port, _ = strconv.Atoi(os.Getenv("COCKROACHDB_PORT"))
	c.User = os.Getenv("COCKROACHDB_USER")
	c.Password = os.Getenv("COCKROACHDB_PASSWORD")
	c.Database = os.Getenv("COCKROACHDB_DATABASE")

	return c
}

func (c Config) URL() string {
	fields := map[string]string{
		"host":     c.Host,
		"port":     fmt.Sprintf("%d", c.Port),
		"user":     c.User,
		"password": c.Password,
		"dbname":   c.Database,
	}

	var dsn string
	for field, value := range fields {
		if value != "" {
			dsn = fmt.Sprintf("%s %s=%s", dsn, field, value)
		}
	}

	return dsn
}

func (c Config) Validate() error {
	if c.User == "" {
		return ErrNoUserInConfig
	}

	if c.Host == "" {
		return ErrNoHostInConfig
	}

	if c.Port == 0 {
		return ErrNoPortInConfig
	}

	return nil
}
