package cmdBacktest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	candlesticksProto "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
)

type AdvanceHandler struct {
	repository vdb.Port
	pubsub     pubsub.Port
	csClient   candlesticksProto.CandlesticksServiceClient
}

func NewAdvanceHandler(repository vdb.Port, ps pubsub.Port, csClient candlesticksProto.CandlesticksServiceClient) AdvanceHandler {
	if repository == nil {
		panic("nil repository")
	}

	if ps == nil {
		panic("nil pubsub")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return AdvanceHandler{
		repository: repository,
		pubsub:     ps,
		csClient:   csClient,
	}
}

func (h AdvanceHandler) Handle(ctx context.Context, backtestId uint) error {
	return h.repository.LockedBacktest(backtestId, func() error {
		// Get backtest info
		bt, err := h.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		// Advance backtest
		finished := bt.Advance()

		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = h.readActualEvents(ctx, bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				log.Println("WARNING: no event detected for", bt.CurrentCsTick.Time)
				bt.SetCurrentTime(bt.EndTime)
			} else if !evts[0].Time.Equal(bt.CurrentCsTick.Time) {
				log.Println("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].Time)
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		evts = append(evts, event.NewStatusEvent(bt.CurrentCsTick.Time, status.Status{
			Finished: finished,
		}))
		h.broadcastEvents(ctx, backtestId, evts)

		if len(evts) > 1 {
			if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
				return fmt.Errorf("cannot update backtest: %w", err)
			}
		}

		return nil
	})
}

func (h AdvanceHandler) readActualEvents(ctx context.Context, bt backtest.Backtest) ([]event.Event, error) {
	evts := make([]event.Event, 0, len(bt.TickSubscribers))
	for _, sub := range bt.TickSubscribers {
		resp, err := h.csClient.ReadCandlesticks(ctx, &candlesticksProto.ReadCandlesticksRequest{
			ExchangeName: sub.ExchangeName,
			PairSymbol:   sub.PairSymbol,
			PeriodSymbol: bt.PeriodBetweenEvents.String(),
			Start:        bt.CurrentCsTick.Time.Format(time.RFC3339),
			End:          bt.EndTime.Format(time.RFC3339),
			Limit:        1,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		if len(resp.Candlesticks) == 0 {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.ExchangeName, sub.PairSymbol, bt.CurrentCsTick.PriceType, resp.Candlesticks[0])
		if err != nil {
			return nil, fmt.Errorf("turning candlestick into event: %w", err)
		}
		evts = append(evts, evt)
	}

	_, evts = event.OnlyKeepEarliestSameTimeEvents(evts, bt.EndTime)
	return evts, nil
}

func (h AdvanceHandler) broadcastEvents(ctx context.Context, backtestId uint, evts []event.Event) {
	var count uint
	for _, evt := range evts {
		if err := h.pubsub.Publish(ctx, backtestId, evt); err != nil {
			log.Println("WARNING: error when publishing event", evt)
			continue
		}

		count++
	}

	if count == 0 {
		log.Println("WARNING: no available events")
	}
}
