package binance

import (
	"fmt"

	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/internal/go/adapters/binance"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
)

const Name = "binance"

type Service struct {
	config binance.Config
	client *client.Client
}

func New() (*Service, error) {
	var c binance.Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading binance config: %w", err)
	}

	return &Service{
		config: c,
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
