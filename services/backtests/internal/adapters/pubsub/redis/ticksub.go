package redis

import (
	"encoding/json"
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/event"
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

func (ps *PriceSubscriber) Channel() <-chan event.Interface {
	ch := make(chan event.Interface)
	go redisChannelToBKChannel(ps.redisPriceSubscriber.Channel(), ch)
	return ch
}

func redisChannelToBKChannel(redisChan <-chan *redis.Message, bkChan chan<- event.Interface) {
	for {
		select {
		case msg, open := <-redisChan:
			if !open {
				close(bkChan)
				return
			}

			var evt event.Event
			if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
				log.Println("error when unmarshaling from redis,", err)
				continue
			}

			switch evt.Type {
			case event.TypeIsTick:
				var tickEvt event.TickEvent
				if err := json.Unmarshal([]byte(msg.Payload), &tickEvt); err != nil {
					log.Println("error when unmarshaling tick from redis,", err)
					continue
				}
				bkChan <- tickEvt
			case event.TypeIsEnd:
				var endEvt event.EndEvent
				if err := json.Unmarshal([]byte(msg.Payload), &endEvt); err != nil {
					log.Println("error when unmarshaling end from redis,", err)
					continue
				}
				bkChan <- endEvt
			}
		}
	}
}

func (ps *PriceSubscriber) Close() error {
	return ps.redisPriceSubscriber.Close()
}
