package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain"
	"golang.org/x/xerrors"
)

type CachedReadCandlesticksPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}

func (req *CachedReadCandlesticksPayload) Process() {
	defaultDuration := req.Period.Duration() * 500
	if req.End == nil {
		if req.Start == nil {
			t := time.Now()
			req.End = &t
		} else {
			t := req.Start.Add(defaultDuration)
			req.End = &t
		}
	}

	if req.Start == nil {
		t := req.End.Add(-defaultDuration)
		req.Start = &t
	}

	t := req.Period.RoundTime(*req.Start)
	req.Start = &t

	t = req.Period.RoundTime(*req.End)
	req.End = &t
}

type CachedReadCandlesticksHandler struct {
	repository db.Port
	services   map[string]exchanges.Port
}

func NewCachedReadCandlesticksHandler(
	repository db.Port,
	services map[string]exchanges.Port,
) CachedReadCandlesticksHandler {
	if repository == nil {
		panic("nil repository")
	}

	if services == nil || len(services) == 0 {
		panic("nil services")
	}

	return CachedReadCandlesticksHandler{
		repository: repository,
		services:   services,
	}
}

func (reh CachedReadCandlesticksHandler) Handle(ctx context.Context, payload CachedReadCandlesticksPayload) (*candlestick.List, error) {
	payload.Process()

	id := candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	}
	cl := candlestick.NewList(id)

	if err := reh.repository.ReadCandlesticks(ctx, cl, *payload.Start, *payload.End, payload.Limit); err != nil {
		return nil, err
	}

	if !domain.AreCsMissing(cl, *payload.Start, *payload.End, payload.Limit) {
		return cl, nil
	}

	downloadEnd := domain.MinimalCandlesticksEndTimeForDownload(cl, *payload.Start, *payload.End)
	if err := reh.download(ctx, cl, *payload.Start, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	if err := reh.upsert(ctx, cl); err != nil {
		return nil, err
	}

	return domain.GetRequestedCandlesticksFromList(cl, *payload.Start, *payload.End, payload.Limit), nil
}

func (reh CachedReadCandlesticksHandler) download(ctx context.Context, cl *candlestick.List, start, end time.Time, limit uint) error {
	exchangeService, exists := reh.services[cl.ExchangeName()]
	if !exists {
		return xerrors.New(fmt.Sprintf("inexistant exchange service for %q", cl.ExchangeName()))
	}

	service, err := exchangeService.Candlesticks(cl.PairSymbol(), cl.Period())
	if err != nil {
		return err
	}

	service.StartTime(start).EndTime(end)
	for {
		ncl, err := service.Do(ctx)
		if err != nil {
			return err
		}

		if err := cl.Merge(*ncl, nil); err != nil {
			return err
		}

		if err := cl.ReplaceUncomplete(*ncl); err != nil {
			return err
		}

		t, _, exists := ncl.Last()
		if !exists || t.Equal(end) || (limit != 0 && cl.Len() >= int(limit)) {
			break
		}

		service.StartTime(t.Add(cl.Period().Duration()))
	}

	return nil
}

func (reh CachedReadCandlesticksHandler) upsert(ctx context.Context, cl *candlestick.List) error {
	start, _, startExists := cl.First()
	end, _, endExists := cl.Last()
	if !startExists || !endExists {
		return nil
	}

	rcl := candlestick.NewList(cl.ID())
	if err := reh.repository.ReadCandlesticks(ctx, rcl, start, end, 0); err != nil {
		return err
	}

	csToInsert := candlestick.NewList(cl.ID())
	csToUpdate := candlestick.NewList(cl.ID())
	if err := cl.Loop(func(ts time.Time, cs candlestick.Candlestick) (bool, error) {
		c, exists := rcl.Get(ts)
		if !exists {
			return false, csToInsert.Set(ts, cs)
		}

		if c.Uncomplete {
			return false, csToUpdate.Set(ts, cs)
		}

		return false, nil
	}); err != nil {
		return err
	}

	if csToInsert.Len() > 0 {
		return reh.repository.CreateCandlesticks(ctx, csToInsert)
	}

	if csToUpdate.Len() > 0 {
		return reh.repository.UpdateCandlesticks(ctx, csToUpdate)
	}

	return nil
}
