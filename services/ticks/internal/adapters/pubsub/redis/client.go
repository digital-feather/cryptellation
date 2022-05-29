package redis

import (
	"context"
	"fmt"

	config "github.com/digital-feather/cryptellation/internal/adapters/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
	"github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
)

const (
	ticksChannelName = "ticks-symbol-%s"
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

func (c *Client) Publish(ctx context.Context, t tick.Tick) error {
	channel := fmt.Sprintf(ticksChannelName, t.PairSymbol)
	return c.client.Publish(ctx, channel, t).Err()
}

func (c *Client) Subscribe(ctx context.Context, symbol string) (pubsub.Subscriber, error) {
	channel := fmt.Sprintf(ticksChannelName, symbol)
	pubsub := c.client.Subscribe(ctx, channel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return newPriceSubscriber(pubsub), nil
}
