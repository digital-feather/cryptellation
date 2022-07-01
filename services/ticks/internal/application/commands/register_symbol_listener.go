package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/internal/domain/tick"
)

type RegisterSymbolListenerHandler struct {
	vdb       vdb.Port
	pubsub    pubsub.Port
	exchanges map[string]exchanges.Port
}

func NewRegisterSymbolListener(ps pubsub.Port, db vdb.Port, exchanges map[string]exchanges.Port) RegisterSymbolListenerHandler {
	if ps == nil {
		panic("nil pubsub")
	}

	if db == nil {
		panic("nil vdb")
	}

	if exchanges == nil {
		panic("nil exchanges clients")
	}

	return RegisterSymbolListenerHandler{
		pubsub:    ps,
		exchanges: exchanges,
		vdb:       db,
	}
}

func (h RegisterSymbolListenerHandler) Handle(ctx context.Context, exchange, pairSymbol string) error {
	count, err := h.vdb.IncrementSymbolListenerCount(ctx, exchange, pairSymbol)
	if err != nil {
		return err
	}

	if count == 1 {
		err := h.launchListener(exchange, pairSymbol)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h RegisterSymbolListenerHandler) launchListener(exchange, pairSymbol string) error {
	exch, exists := h.exchanges[exchange]
	if !exists {
		return fmt.Errorf("exchange %q doesn't exists", exchange)
	}

	ticksChan, stopChan, err := exch.ListenSymbol(pairSymbol)
	if err != nil {
		return fmt.Errorf("listening to symbol: %w", err)
	}

	go newListener(newListenerPayload{
		exchange:   exchange,
		pairSymbol: pairSymbol,
		db:         h.vdb,
		ps:         h.pubsub,
		ticksChan:  ticksChan,
		stopChan:   stopChan,
	})

	return nil
}

type newListenerPayload struct {
	exchange   string
	pairSymbol string
	db         vdb.Port
	ps         pubsub.Port
	ticksChan  chan tick.Tick
	stopChan   chan struct{}
}

const checkInterval = 10 * time.Second

func newListener(payload newListenerPayload) {
	ctx := context.Background()
	nextCheckTime := time.Now().Add(checkInterval)
	lastPrice := 0.0

	for {
		t := <-payload.ticksChan
		if lastPrice == t.Price {
			continue
		}
		lastPrice = t.Price

		err := payload.ps.Publish(ctx, t)
		if err != nil {
			log.Println("Publish error:", err)
			continue
		}

		if nextCheckTime.Before(time.Now()) {
			count, err := payload.db.GetSymbolListenerCount(ctx, payload.exchange, payload.pairSymbol)
			if err != nil {
				log.Println(err)
				continue
			}

			if count <= 0 {
				log.Println("Interrupting", payload.exchange, payload.pairSymbol, "listener")
				payload.stopChan <- struct{}{}
				break
			}

			nextCheckTime = time.Now().Add(checkInterval)
		}
	}
}
