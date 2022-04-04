package binance

import (
	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/internal/adapters/binance"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/period"
	"golang.org/x/xerrors"
)

const Name = "binance"

type Service struct {
	config binance.Config
	client *client.Client
}

func New() (*Service, error) {
	var c binance.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading binance config: %w", err)
	}

	return &Service{
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, nil
}

func (s *Service) Candlesticks(pairSymbol string, per period.Symbol) (exchanges.CandlesticksService, error) {
	service := s.client.NewKlinesService()
	service.Symbol(BinanceSymbol(pairSymbol))

	binanceInterval, err := PeriodToInterval(per)
	if err != nil {
		return nil, err
	}
	service.Interval(binanceInterval)

	return &CandlestickService{
		service:    service,
		pairSymbol: pairSymbol,
		period:     per,
	}, nil
}
