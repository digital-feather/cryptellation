package redis

import (
	"context"
	"fmt"

	config "github.com/digital-feather/cryptellation/internal/adapters/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
	"github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
)

const (
	pricesChannelName = "backtest-%d-ticks"
)

type Client struct {
	client *redis.Client
}

func New() (*Client, error) {
	var c config.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading redis config: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	// TODO Check connection

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Publish(ctx context.Context, backtestID uint, evt event.Interface) error {
	channel := fmt.Sprintf(pricesChannelName, backtestID)
	return c.client.Publish(ctx, channel, evt).Err()
}

func (c *Client) Subscribe(ctx context.Context, backtestID uint) (pubsub.Subscriber, error) {
	channel := fmt.Sprintf(pricesChannelName, backtestID)
	pubsub := c.client.Subscribe(ctx, channel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return newPriceSubscriber(pubsub), nil
}
