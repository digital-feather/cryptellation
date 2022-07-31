package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	config "github.com/digital-feather/cryptellation/internal/go/adapters/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
	"github.com/go-redis/redis/v8"
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
		return nil, fmt.Errorf("loading redis config: %w", err)
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

	content, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream:       channel,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"content": content,
		},
	}).Err()
}

func (c *Client) Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick)
	go c.redisToChannelEvents(ctx, symbol, ch)
	time.Sleep(10 * time.Millisecond) // Wait 1 millisecond to avoid missing messages
	return ch, nil
}

func (c *Client) redisToChannelEvents(ctx context.Context, symbol string, ch chan tick.Tick) {
	id := "$"
	channel := fmt.Sprintf(ticksChannelName, symbol)
	for {
		// Reading next message
		cmd := c.client.XRead(context.TODO(), &redis.XReadArgs{
			Streams: []string{channel, id},
			Count:   1,
			Block:   time.Minute,
		})
		err := cmd.Err()
		if err != nil {
			continue
		}

		// Exit if context is done or continue
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
		}

		// Setting next ID
		msg := cmd.Val()[0].Messages[0]
		id = msg.ID

		// Passing the tick
		t, err := msgToTick(msg.Values)
		if err != nil {
			log.Println("malformed tick:", err)
			continue
		}
		ch <- t
	}
}

func msgToTick(m map[string]interface{}) (tick.Tick, error) {
	var t tick.Tick

	content, ok := m["content"]
	if !ok {
		return tick.Tick{}, fmt.Errorf("there is no content in the message: %v", m)
	}

	contentStr, ok := content.(string)
	if !ok {
		return tick.Tick{}, fmt.Errorf("'content' field is no string in message")
	}

	err := json.Unmarshal([]byte(contentStr), &t)
	return t, err
}
