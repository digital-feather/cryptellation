package cmdBacktest

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/internal/genproto/candlesticks"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/order"
	"golang.org/x/xerrors"
)

type CreateOrderHandler struct {
	repository vdb.Port
	csClient   candlesticks.CandlesticksServiceClient
}

func NewCreateOrderHandler(repository vdb.Port, csClient candlesticks.CandlesticksServiceClient) CreateOrderHandler {
	if repository == nil {
		panic("nil repository")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return CreateOrderHandler{
		repository: repository,
		csClient:   csClient,
	}
}

type CreateOrderPayload struct {
	BacktestId   uint
	Type         order.Type
	ExchangeName string
	PairSymbol   string
	Side         order.Side
	Quantity     float64
}

func (p CreateOrderPayload) ToOrder() order.Order {
	return order.Order{
		Type:         p.Type,
		ExchangeName: p.ExchangeName,
		PairSymbol:   p.PairSymbol,
		Side:         p.Side,
		Quantity:     p.Quantity,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, payload CreateOrderPayload) error {
	ord := payload.ToOrder()
	if err := ord.Validate(); err != nil {
		return xerrors.Errorf("invalid order: %w", err)
	}

	return h.repository.LockedBacktest(payload.BacktestId, func() error {
		bt, err := h.repository.ReadBacktest(ctx, payload.BacktestId)
		if err != nil {
			return xerrors.Errorf("cannot get backtest: %w", err)
		}

		resp, err := h.csClient.ReadCandlesticks(ctx, &candlesticks.ReadCandlesticksRequest{
			ExchangeName: payload.ExchangeName,
			PairSymbol:   payload.PairSymbol,
			PeriodSymbol: bt.PeriodBetweenEvents.String(),
			Start:        bt.CurrentCsTick.Time.Format(time.RFC3339),
			End:          bt.CurrentCsTick.Time.Format(time.RFC3339),
			Limit:        0,
		})
		if err != nil {
			return xerrors.Errorf("could not get candlesticks from service: %w", err)
		} else if len(resp.Candlesticks) == 0 {
			return backtest.ErrNoDataForOrderValidation
		}

		if err := bt.AddOrder(ord, *resp.Candlesticks[0]); err != nil {
			return err
		}

		if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
			return xerrors.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
