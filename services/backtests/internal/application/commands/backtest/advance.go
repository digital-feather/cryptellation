package cmdBacktest

import (
	"context"
	"log"
	"time"

	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/cryptellation/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/cryptellation/cryptellation/services/backtests/internal/domain/event"
	"golang.org/x/xerrors"
)

type AdvanceHandler struct {
	repository vdb.Port
	pubsub     pubsub.Port
	csClient   candlesticks.CandlesticksServiceClient
}

func NewAdvanceHandler(repository vdb.Port, ps pubsub.Port, csClient candlesticks.CandlesticksServiceClient) AdvanceHandler {
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

func (h AdvanceHandler) Handle(ctx context.Context, backtestId uint) (finished bool, err error) {
	err = h.repository.LockedBacktest(backtestId, func() error {
		bt, err := h.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return xerrors.Errorf("cannot get backtest: %w", err)
		}

		bt.Advance()
		if bt.Done() {
			finished = true
			return nil
		}

		evts, err := h.readActualEvents(ctx, backtestId)
		if err != nil {
			return xerrors.Errorf("cannot read actual events: %w", err)
		}
		h.broadcastEvents(ctx, backtestId, evts)

		if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
			return xerrors.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})

	return finished, err
}

func (h AdvanceHandler) readActualEvents(ctx context.Context, backtestId uint) ([]event.Interface, error) {
	bt, err := h.repository.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	evts := make([]event.Interface, len(bt.TickSubscribers))
	for i, sub := range bt.TickSubscribers {
		resp, err := h.csClient.ReadCandlesticks(ctx, &candlesticks.ReadCandlesticksRequest{
			ExchangeName: sub.ExchangeName,
			PairSymbol:   sub.PairSymbol,
			PeriodSymbol: bt.PeriodBetweenEvents.String(),
			Start:        bt.CurrentCsTick.Time.Format(time.RFC3339),
			End:          bt.EndTime.Format(time.RFC3339),
			Limit:        1,
		})
		if err != nil {
			return nil, xerrors.Errorf("could not get candlesticks from service: %w", err)
		}

		if len(resp.Candlesticks) == 0 {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.ExchangeName, sub.PairSymbol, bt.CurrentCsTick.PriceType, *resp.Candlesticks[0])
		if err != nil {
			return nil, xerrors.Errorf("turning candlestick into event: %w", err)
		}
		evts[i] = evt
	}

	t, evts := event.OnlyKeepEarliestSameTimeEvents(evts)
	if len(evts) == 0 {
		log.Println("WARNING: no event detected for", bt.CurrentCsTick.Time)
		bt.SetCurrentTime(bt.EndTime)
	} else if !evts[0].GetTime().Equal(bt.CurrentCsTick.Time) {
		log.Println("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].GetTime())
		bt.SetCurrentTime(evts[0].GetTime())
	}

	return append(evts, event.NewEndEvent(t)), nil
}

func (h AdvanceHandler) broadcastEvents(ctx context.Context, backtestId uint, evts []event.Interface) {
	var count uint
	for _, evt := range evts {
		if evt == nil {
			continue
		}

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
