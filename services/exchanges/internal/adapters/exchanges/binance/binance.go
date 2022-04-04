package binance

import (
	"context"
	"fmt"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/internal/adapters/binance"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"golang.org/x/xerrors"
)

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

func (ps *Service) Infos(ctx context.Context) (exchange.Exchange, error) {
	exchangeInfos, err := ps.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchange.Exchange{}, err
	}

	pairSymbols := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairSymbols[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := exchanges.Binance
	exch.PairsSymbols = pairSymbols
	exch.LastSyncTime = time.Now()

	return exch, nil
}
