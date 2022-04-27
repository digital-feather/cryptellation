package redis

import (
	"encoding/json"
	"log"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
	"github.com/go-redis/redis/v8"
)

type PriceSubscriber struct {
	redisPriceSubscriber *redis.PubSub
}

func newPriceSubscriber(pubsub *redis.PubSub) pubsub.Subscriber {
	return &PriceSubscriber{
		redisPriceSubscriber: pubsub,
	}
}

func (ps *PriceSubscriber) Channel() <-chan tick.Tick {
	ch := make(chan tick.Tick)
	go redisChannelToBKChannel(ps.redisPriceSubscriber.Channel(), ch)
	return ch
}

func redisChannelToBKChannel(redisChan <-chan *redis.Message, bkChan chan<- tick.Tick) {
	for {
		select {
		case msg, open := <-redisChan:
			if !open {
				close(bkChan)
				return
			}

			var t tick.Tick
			if err := json.Unmarshal([]byte(msg.Payload), &t); err != nil {
				log.Println("error when unmarshaling from redis,", err)
				continue
			}

			bkChan <- t
		}
	}
}

func (ps *PriceSubscriber) Close() error {
	return ps.redisPriceSubscriber.Close()
}
