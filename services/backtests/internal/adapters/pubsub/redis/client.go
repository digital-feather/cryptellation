package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	config "github.com/digital-feather/cryptellation/internal/adapters/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/tick"
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
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream:       channel,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"type":    evt.GetType(),
			"time":    evt.GetTime(),
			"content": evt.GetContent(),
		},
	}).Err()
}

func (c *Client) Subscribe(ctx context.Context, backtestID uint) (<-chan event.Interface, error) {
	ch := make(chan event.Interface)
	go c.redisToChannelEvents(ctx, backtestID, ch)
	time.Sleep(time.Millisecond) // Wait 1 millisecond to avoid missing messages
	return ch, nil
}

func (c *Client) redisToChannelEvents(ctx context.Context, backtestID uint, ch chan event.Interface) {
	id := "$"
	channel := fmt.Sprintf(pricesChannelName, backtestID)
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

		// Passing the event
		switch event.Type(fmt.Sprint(msg.Values["type"])) {
		case event.TypeIsTick:
			tickEvent, err := eventToTick(msg.Values)
			if err != nil {
				log.Println("error when unmarshaling tick from redis,", err)
				continue
			}
			ch <- tickEvent
		case event.TypeIsEnd:
			endEvent, err := eventToEnd(msg.Values)
			if err != nil {
				log.Println("error when unmarshaling end from redis,", err)
				continue
			}
			ch <- endEvent
		default:
			// TODO: handle
			fmt.Println("Unknown type:", msg.Values["type"])
		}
	}
}

func eventToTick(values map[string]interface{}) (event.TickEvent, error) {
	var ti tick.Tick
	content := fmt.Sprint(values["content"])
	if err := json.Unmarshal([]byte(content), &ti); err != nil {
		return event.TickEvent{}, err
	}

	tString := fmt.Sprint(values["time"])
	t, err := time.Parse(time.RFC3339, tString)
	if err != nil {
		return event.TickEvent{}, err
	}

	return event.NewTickEvent(t, ti), nil
}

func eventToEnd(values map[string]interface{}) (event.EndEvent, error) {
	tString := fmt.Sprint(values["time"])
	t, err := time.Parse(time.RFC3339, tString)
	if err != nil {
		return event.EndEvent{}, err
	}

	return event.NewEndEvent(t), nil
}
